package queues

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/redis/rueidis"
	"github.com/segmentio/ksuid"
)

const valueSeparator = "$$"

var (
	// ErrFailedToAcquireLock is returned when the lock cannot be acquired
	ErrFailedToAcquireLock = errors.New("failed to acquire lock")

	// ErrFailedToParseValue is returned when the redis value cannot be parsed
	ErrFailedToParseValue = errors.New("failed to parse value")

	// ErrFailedToEncodeQueueEntry is returned when the queue entry cannot be marshaled
	ErrFailedToEncodeQueueEntry = errors.New("failed to encode queue entry")

	// ErrFailedToDecodeQueueEntry is returned when the queue entry cannot be marshaled
	ErrFailedToDecodeQueueEntry = errors.New("failed to decode queue entry")

	// ErrFailedExecuteCommand is returned when the redis command cannot be executed
	ErrFailedExecuteCommand = errors.New("failed to execute command")
)

// RedisQueue is the Matchmaker implementation using Redis
type RedisQueue struct {
	config          domain.QueueConfig
	client          rueidis.Client
	lock            ports.Lock
	keyPrefix       string
	allKeys         []string
	bracketInterval int
	codec           ports.Codec
}

// NewRedisQueue creates a new RedisQueue
func NewRedisQueue(c rueidis.Client, queueConfig domain.QueueConfig, codec ports.Codec, lock ports.Lock) *RedisQueue {
	allKeys := allKeysSetup(queueConfig)
	bracketInterval := (queueConfig.MaxRanking - queueConfig.MinRanking) / queueConfig.NrBrackets
	fmt.Println("bracketInterval", bracketInterval)
	return &RedisQueue{
		client:          c,
		config:          queueConfig,
		lock:            lock,
		keyPrefix:       "ranking::" + queueConfig.Name,
		allKeys:         allKeys,
		bracketInterval: bracketInterval,
		codec:           codec,
	}
}

// Enqueue adds a player to the queue
func (q *RedisQueue) Enqueue(ctx context.Context, qEntry domain.QueueEntry) error {
	bytes, err := q.codec.Encode(qEntry)
	if err != nil {
		return errors.Join(ErrFailedToEncodeQueueEntry, err)
	}
	bracket := qEntry.Ranking / q.bracketInterval
	key := q.keyName(bracket)
	return q.internalEnqueue(ctx, key, string(bytes))
}

func (q *RedisQueue) internalEnqueue(ctx context.Context, key string, value string) error {
	v := fmt.Sprintf("%s%s%s", value, valueSeparator, key)
	cmdPush := q.client.B().Rpush().Key(key).Element(v).Build()
	err := q.client.Do(ctx, cmdPush).Error()
	if err != nil {
		return errors.Join(ErrFailedExecuteCommand, fmt.Errorf("when enqueueing"), err)
	}
	return nil
}

// Make finds a match
func (q *RedisQueue) Make(ctx context.Context) (matches []domain.MatchResult, err error) {

	// Acquire a lock
	ctxLock, cancel, err := q.lock.Acquire(ctx, q.config.Name+":lock")
	if err != nil {
		return nil, errors.Join(ErrFailedToAcquireLock, err)
	}
	defer cancel()

	err = nil
	remainingTickets := make([]string, 0)
	matches = make([]domain.MatchResult, 0)

	defer func() {
		count := len(remainingTickets)
		if count > 0 {
			for _, v := range remainingTickets {
				kv := strings.Split(v, valueSeparator)
				if len(kv) != 2 {
					continue
				}
				value := kv[0]
				key := kv[1]
				err = q.internalEnqueue(ctxLock, key, value)
			}
		}
	}()

	// iterate over chunks of keys to given how redis work with LMPop we want to traverse all the keys
	// here we are using a chunk of 3 keys
	chunkSize := 3
	iteration := 0

	for {
		chunkKeys := q.allKeys[iteration*chunkSize : (iteration+1)*chunkSize]
		chunkKeysLen := len(chunkKeys)
		if iteration*chunkSize >= len(q.allKeys) {
			break
		}
		for {

			// using LMPop to get the first N players in the queue (from all the chunk keys)
			// this allows to always start from lower ranking to the max ranking
			cmd := q.client.B().Lmpop().Numkeys(int64(chunkKeysLen)).Key(chunkKeys...).Left().Count(int64(q.config.MaxPlayers)).Build()
			result, err := q.client.Do(ctxLock, cmd).AsMap()
			if err != nil {
				if rueidis.IsRedisNil(err) {
					iteration++
					break
				}
				return nil, errors.Join(ErrFailedExecuteCommand, fmt.Errorf("when getting the keys"), err)
			}

			for _, value := range result {
				v, err := value.AsStrSlice()

				if err != nil {
					return nil, errors.Join(ErrFailedToParseValue, err)
				}
				v = append(remainingTickets, v...)
				remainingTickets = nil
				if len(v) < q.config.MaxPlayers {
					remainingTickets = append(remainingTickets, v...)
					break
				}
				for i := 0; i < len(v); i = i + q.config.MaxPlayers {
					if i <= len(v)-q.config.MaxPlayers {
						entries := make([]domain.QueueEntry, 0)
						for j := 0; j < q.config.MaxPlayers; j++ {
							kv := strings.Split(v[i+j], valueSeparator)
							if len(kv) != 2 {
								return nil, ErrFailedToParseValue
							}
							var qe domain.QueueEntry
							err = q.codec.Decode([]byte(kv[0]), &qe)
							if err != nil {
								return nil, errors.Join(ErrFailedToDecodeQueueEntry, err)
							}
							entries = append(entries, qe)
						}
						matches = append(matches, domain.MatchResult{
							Match: domain.Match{
								ID: ksuid.New().String(),
							},
							Entries: entries,
						})
					} else {
						remainingTickets = v[i:]
					}
				}
			}
		}
	}
	return
}

// Name returns the name of the queue
func (q *RedisQueue) Name() string {
	return q.config.Name
}

// keyName returns the key name for the queue based on the bracket
func (q *RedisQueue) keyName(bracket int) string {
	return q.keyPrefix + "::" + strconv.FormatInt(int64(bracket), 10)
}

// allKeysSetup returns all the keys for the queue based on the number of brackets
func allKeysSetup(queueConfig domain.QueueConfig) []string {
	allKeys := make([]string, 0)

	for i := 0; i < queueConfig.NrBrackets; i++ {
		bracket := fmt.Sprintf("ranking::queue::%s::%d", queueConfig.Name, i)
		allKeys = append(allKeys, bracket)
	}
	slices.Sort(allKeys)
	allKeys = slices.Compact(allKeys)
	return allKeys
}
