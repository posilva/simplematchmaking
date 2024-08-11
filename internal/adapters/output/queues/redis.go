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
	prefix := "ranking::" + queueConfig.Name
	allKeys := allKeysSetup(prefix, queueConfig)
	bracketInterval := (queueConfig.MaxRanking - queueConfig.MinRanking) + 1/queueConfig.NrBrackets
	return &RedisQueue{
		client:          c,
		config:          queueConfig,
		lock:            lock,
		keyPrefix:       prefix,
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
	key := keyName(q.keyPrefix, bracket)
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
	err = nil
	remainingTickets := make([]string, 0)
	matches = make([]domain.MatchResult, 0)
	defer func() {
		cancel()
	}()

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
	chunkSize := q.config.MakeIterations
	iteration := 0
	debugCounter := 0
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
				// convert to a list of strings
				current, err := value.AsStrSlice()
				if err != nil {
					return nil, errors.Join(ErrFailedToParseValue, err)
				}

				debugCounter++
				// add the remaining tickets to the current list
				current = append(remainingTickets, current...)
				lenCurrent := len(current)
				// clean up remaining tickets
				remainingTickets = nil
				// if we have less than the max players to handle then we add to the remaining tickets to be process later
				if lenCurrent < q.config.MaxPlayers {
					remainingTickets = append(remainingTickets, current...)
					break
				}
				// pick first N players (MaxPlayers) to create a match
				for i := 0; i < lenCurrent; i = i + q.config.MaxPlayers {
					if i <= lenCurrent-q.config.MaxPlayers {
						entries := make([]domain.QueueEntry, 0)
						tickets := make([]string, 0)
						for j := 0; j < q.config.MaxPlayers; j++ {
							kv := strings.Split(current[i+j], valueSeparator)
							if len(kv) != 2 {
								return nil, ErrFailedToParseValue
							}
							var qe domain.QueueEntry
							err = q.codec.Decode([]byte(kv[0]), &qe)
							if err != nil {
								return nil, errors.Join(ErrFailedToDecodeQueueEntry, err)
							}
							entries = append(entries, qe)
							tickets = append(tickets, qe.TicketID)
						}
						matches = append(matches, domain.MatchResult{
							Match: domain.Match{
								ID:        ksuid.New().String(),
								TicketIDs: tickets,
							},
							Entries: entries,
						})
					} else {
						remainingTickets = current[i:]
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
func keyName(prefix string, bracket int) string {
	return prefix + "::" + strconv.FormatInt(int64(bracket), 10)
}

// allKeysSetup returns all the keys for the queue based on the number of brackets
func allKeysSetup(prefix string, queueConfig domain.QueueConfig) []string {
	allKeys := make([]string, 0)

	for i := 0; i < queueConfig.NrBrackets; i++ {
		bracket := keyName(prefix, i)
		allKeys = append(allKeys, bracket)
	}
	allKeys = slices.Compact(allKeys)
	return allKeys
}
