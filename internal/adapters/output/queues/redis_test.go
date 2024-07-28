package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/rand"
)

func TestRedisQueue_Enqueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	entry := domain.QueueEntry{
		TicketID: testutil.NewID(),
		PlayerID: testutil.NewID(),
		Ranking:  1,
	}
	bytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}
	name := "test"
	client := mock.NewClient(ctrl)
	v := fmt.Sprintf("%s$$%s", string(bytes), "test::0")
	client.EXPECT().Do(gomock.Any(), mock.Match(
		"RPUSH", "test::0", v)).Return(mock.Result(mock.RedisInt64(1)))

	q := &RedisQueue{
		bracketInterval: 10,
		keyPrefix:       name,
		client:          client,
		config: domain.QueueConfig{
			Name:           name,
			MaxPlayers:     2,
			NrBrackets:     10,
			MinRanking:     1,
			MaxRanking:     100,
			MakeIterations: 3,
		},
	}
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
	bytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}
	name := "test"
	client := mock.NewClient(ctrl)
	v := fmt.Sprintf("%s$$%s", string(bytes), "test::0")
	client.EXPECT().Do(gomock.Any(), mock.Match(
		"RPUSH", "test::0", v)).Return(mock.ErrorResult(fmt.Errorf("error")))

	q := &RedisQueue{
		bracketInterval: 10,
		keyPrefix:       name,
		client:          client,
		config: domain.QueueConfig{
			Name:           name,
			MaxPlayers:     2,
			NrBrackets:     10,
			MinRanking:     1,
			MaxRanking:     100,
			MakeIterations: 3,
		},
	}
	err = q.Enqueue(context.Background(), entry)
	require.ErrorContains(t, err, ErrFailedExecuteCommand.Error())

}
func TestRedisQueue_Make(t *testing.T) {
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
						//t.Logf("match between: %s ", v[i:i+int(count)])
					} else {
						remainingPlayers = v[i:]
					}
				}
			}
		}
	}

	require.Equal(t, len(allMatched), maxPlayers/int(count))
}
