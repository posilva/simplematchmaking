package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/domain/codecs"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/posilva/simplematchmaking/internal/core/ports/mocks"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/rand"
)

var (
	name        = "test"
	queueConfig = domain.QueueConfig{
		Name:           name,
		MaxPlayers:     2,
		NrBrackets:     10,
		MinRanking:     1,
		MaxRanking:     100,
		MakeIterations: 3,
	}
)

func TestRedisQueue_Enqueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	codec := codecs.NewJSONCodec()
	entry := domain.QueueEntry{
		TicketID: testutil.NewID(),
		PlayerID: testutil.NewID(),
		Ranking:  1,
	}
	bytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}
	lock := mocks.NewMockLock(ctrl)
	client := mock.NewClient(ctrl)
	v := fmt.Sprintf("%s$$%s", string(bytes), "ranking::test::0")
	client.EXPECT().Do(gomock.Any(), mock.Match(
		"RPUSH", "ranking::test::0", v)).Return(mock.Result(mock.RedisInt64(1)))

	q := NewRedisQueue(client, queueConfig, codec, lock)
	err = q.Enqueue(context.Background(), entry)
	require.NoError(t, err)

}
func TestRedisQueue_Enqueue_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entry := domain.QueueEntry{
		TicketID: testutil.NewID(),
		PlayerID: testutil.NewID(),
		Ranking:  1,
	}
	codec := codecs.NewJSONCodec()

	bytes, err := codec.Encode(entry)
	if err != nil {
		t.Fatal(err)
	}
	lock := mocks.NewMockLock(ctrl)
	client := mock.NewClient(ctrl)
	v := fmt.Sprintf("%s$$%s", string(bytes), "ranking::test::0")
	client.EXPECT().Do(gomock.Any(), mock.Match(
		"RPUSH", "ranking::test::0", v)).Return(mock.ErrorResult(fmt.Errorf("error")))

	q := NewRedisQueue(client, queueConfig, codec, lock)
	err = q.Enqueue(context.Background(), entry)
	require.ErrorContains(t, err, ErrFailedExecuteCommand.Error())

}
func TestRedisQueue_Enqueue_Error_Encode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entry := domain.QueueEntry{
		TicketID: testutil.NewID(),
		PlayerID: testutil.NewID(),
		Ranking:  1,
	}
	codecMock := mocks.NewMockCodec(ctrl)
	client := mock.NewClient(ctrl)
	lock := mocks.NewMockLock(ctrl)
	codecMock.EXPECT().Encode(gomock.Any()).Return(nil, fmt.Errorf("wrong"))

	q := NewRedisQueue(client, queueConfig, codecMock, lock)
	err := q.Enqueue(context.Background(), entry)
	require.ErrorContains(t, err, ErrFailedToEncodeQueueEntry.Error())

}

func TestRedisQueue_Make(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lockMock := mocks.NewMockLock(ctrl)
	clientMock := mock.NewClient(ctrl)

	codec := codecs.NewJSONCodec()

	lockMock.EXPECT().Acquire(gomock.Any(), "test:lock").
		Return(context.Background(), func() {
			fmt.Println("lock was canceled by the test")
		}, nil)

	entry, entryS := queueEntry(t, codec)
	entry2, entry2S := queueEntry(t, codec)

	listResult := mock.Result(
		mock.RedisMap(
			mapResult(entryS, entry2S),
		))

	nilResult := mock.Result(
		mock.RedisNil(),
	)

	ct := 0
	clientMock.EXPECT().
		Do(gomock.Any(),
			mock.MatchFn(
				func(cmd []string) bool {
					return true
				}, "testing description of matcher fn",
			)).AnyTimes().
		DoAndReturn(func(ctx context.Context, cmd interface{}) rueidis.RedisResult {
			ct++
			if ct > 1 {
				return nilResult
			}

			return listResult
		})

	q := NewRedisQueue(clientMock, queueConfig, codec, lockMock)
	matches, err := q.Make(context.Background())
	require.NoError(t, err)
	require.Len(t, matches, 1)
	require.Len(t, matches[0].Entries, 2)

	// it's time to decode tickets received
	qe1 := matches[0].Entries[0]
	qe2 := matches[0].Entries[1]

	require.Equal(t, entry, qe1)
	require.Equal(t, entry2, qe2)

}

func TestRedisQueue_Make_LockError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lockMock := mocks.NewMockLock(ctrl)
	clientMock := mock.NewClient(ctrl)

	codec := codecs.NewJSONCodec()

	lockMock.EXPECT().Acquire(gomock.Any(), gomock.Any()).
		Return(context.Background(), func() {
			fmt.Println("lock was canceled by the test")
		}, fmt.Errorf("lock error"))

	q := NewRedisQueue(clientMock, queueConfig, codec, lockMock)
	_, err := q.Make(context.Background())
	require.ErrorContains(t, err, ErrFailedToAcquireLock.Error())

}

func TestRedisQueue_Make_ParseErrorStrSlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lockMock := mocks.NewMockLock(ctrl)
	clientMock := mock.NewClient(ctrl)

	codec := codecs.NewJSONCodec()

	lockMock.EXPECT().Acquire(gomock.Any(), "test:lock").
		Return(context.Background(), func() {
			fmt.Println("lock was canceled by the test")
		}, nil)

	_, entryS := queueEntry(t, codec)
	_, entry2S := queueEntry(t, codec)

	listResult := mock.Result(
		mock.RedisMap(
			mapResultNotStrSlice(entryS, entry2S),
		))

	nilResult := mock.Result(
		mock.RedisNil(),
	)

	ct := 0
	clientMock.EXPECT().
		Do(gomock.Any(),
			mock.MatchFn(
				func(cmd []string) bool {
					return true
				}, "testing description of matcher fn",
			)).AnyTimes().
		DoAndReturn(func(ctx context.Context, cmd interface{}) rueidis.RedisResult {
			ct++
			if ct > 1 {
				return nilResult
			}

			return listResult
		})

	q := NewRedisQueue(clientMock, queueConfig, codec, lockMock)
	_, err := q.Make(context.Background())
	require.ErrorContains(t, err, ErrFailedToParseValue.Error())

}

func TestRedisQueue_Make_ParseErrorParseEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lockMock := mocks.NewMockLock(ctrl)
	clientMock := mock.NewClient(ctrl)

	codec := codecs.NewJSONCodec()

	lockMock.EXPECT().Acquire(gomock.Any(), "test:lock").
		Return(context.Background(), func() {
			fmt.Println("lock was canceled by the test")
		}, nil)

	_, entryS := queueEntry(t, codec)
	_, entry2S := queueEntry(t, codec)

	listResult := mock.Result(
		mock.RedisMap(
			mapResultNotParse(entryS, entry2S),
		))

	nilResult := mock.Result(
		mock.RedisNil(),
	)

	ct := 0
	clientMock.EXPECT().
		Do(gomock.Any(),
			mock.MatchFn(
				func(cmd []string) bool {
					return true
				}, "testing description of matcher fn",
			)).AnyTimes().
		DoAndReturn(func(ctx context.Context, cmd interface{}) rueidis.RedisResult {
			ct++
			if ct > 1 {
				return nilResult
			}

			return listResult
		})

	q := NewRedisQueue(clientMock, queueConfig, codec, lockMock)
	_, err := q.Make(context.Background())
	require.ErrorContains(t, err, ErrFailedToParseValue.Error())

}
func TestRedisQueue_Make_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lockMock := mocks.NewMockLock(ctrl)
	clientMock := mock.NewClient(ctrl)

	codec := codecs.NewJSONCodec()

	lockMock.EXPECT().Acquire(gomock.Any(), "test:lock").
		Return(context.Background(), func() {
			fmt.Println("lock was canceled by the test")
		}, nil)

	clientMock.EXPECT().
		Do(gomock.Any(),
			mock.MatchFn(
				func(cmd []string) bool {
					t.Log("redis command:", cmd)
					return cmd[0] == "LMPOP"
				}, "testing description of matcher fn",
			)).AnyTimes().Return(mock.ErrorResult(fmt.Errorf("error")))

	q := NewRedisQueue(clientMock, queueConfig, codec, lockMock)
	_, err := q.Make(context.Background())
	require.ErrorContains(t, err, ErrFailedExecuteCommand.Error())
}

func TestRedisQueue_MakeExperiment(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping test in short mode as this is to run in a local redis instance")
	}
	count := int64(8)
	// player with ranking n
	maxPlayers := 1000
	nrBrackets := 1000
	maxRanking := 9999
	minRanking := 1

	bracketInterval := maxRanking / nrBrackets

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := client.Do(context.Background(), client.B().Flushall().Build()).Error()
		if err != nil {
			t.Fatal(err)
		}
		client.Close()
	}()

	keyUnique := testutil.NewID()

	listOfKeys := make([]string, 0)
	listOfUIDs := make([]string, 0)

	// generate a list of maxPlayers
	for i := 0; i < maxPlayers; i++ {
		ranking := rand.Intn(maxRanking-minRanking) + minRanking
		slot := ranking / bracketInterval
		key := fmt.Sprintf("test:list:%s:%d", keyUnique, slot)
		uid := testutil.NewID()
		listOfKeys = append(listOfKeys, key)
		listOfUIDs = append(listOfUIDs, uid)
		cmdPush := client.B().Rpush().Key(key).Element(uid).Build()
		err = client.Do(context.Background(), cmdPush).Error()
		if err != nil {
			t.Fatal(err)
		}
	}

	slices.Sort(listOfKeys)
	remainingPlayers := make([]string, 0)
	allMatched := make([]string, 0)

	listOfKeys = slices.Compact(listOfKeys)

	iterationSize := 3
	iteration := 0
	for {
		keys := listOfKeys[iteration*iterationSize : (iteration+1)*iterationSize]
		lAllKeys := len(keys)
		if iteration*iterationSize >= len(listOfKeys) {
			t.Logf("remaining players: %v", remainingPlayers)
			break
		}
		for {
			cmd := client.B().Lmpop().Numkeys(int64(lAllKeys)).Key(keys...).Left().Count(count).Build()
			result, err := client.Do(context.Background(), cmd).AsMap()
			if err != nil {
				if rueidis.IsRedisNil(err) {
					iteration++
					break
				}
				t.Fatal(err)
			}

			for _, value := range result {
				v, err := value.AsStrSlice()
				if err != nil {
					t.Fatal(err)
				}
				v = append(remainingPlayers, v...)
				remainingPlayers = nil
				if len(v) < int(count) {
					remainingPlayers = append(remainingPlayers, v...)
					break
				}
				for i := 0; i < len(v); i = i + int(count) {
					if i <= len(v)-int(count) {
						allMatched = append(allMatched, v[i]+":"+v[i+1])
					} else {
						remainingPlayers = v[i:]
					}
				}
			}
		}
	}

	require.Equal(t, len(allMatched), maxPlayers/int(count))
}

func queueEntry(t *testing.T, codec ports.Codec) (domain.QueueEntry, string) {
	entry := domain.QueueEntry{
		TicketID: testutil.NewID(),
		PlayerID: testutil.NewID(),
		Ranking:  1,
	}
	b, err := codec.Encode(entry)
	require.NoError(t, err)
	return entry, string(b)
}

func mapResult(entries ...string) map[string]rueidis.RedisMessage {
	kvReturn := make(map[string]rueidis.RedisMessage)
	keyName := "ranking::test::0"
	messages := make([]rueidis.RedisMessage, 0)

	for _, v := range entries {
		messages = append(messages, mock.RedisString(v+"$$"+keyName))
	}
	kvReturn[keyName] = mock.RedisArray(messages...)
	return kvReturn
}
func mapResultNotParse(entries ...string) map[string]rueidis.RedisMessage {
	kvReturn := make(map[string]rueidis.RedisMessage)
	keyName := "ranking::test::0"
	messages := make([]rueidis.RedisMessage, 0)

	for _ = range entries {
		messages = append(messages, mock.RedisString("empty"))
	}
	kvReturn[keyName] = mock.RedisArray(messages...)
	return kvReturn
}
func mapResultNotStrSlice(entries ...string) map[string]rueidis.RedisMessage {
	kvReturn := make(map[string]rueidis.RedisMessage)
	keyName := "ranking::test::0"
	kvReturn[keyName] = mock.RedisFloat64(1)
	return kvReturn
}
