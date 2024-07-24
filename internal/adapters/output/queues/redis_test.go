package queues

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"go.uber.org/mock/gomock"
)

func TestRedisQueue_AddPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock.NewClient(ctrl)
	client.EXPECT().Do(gomock.Any(), mock.Match(
		"ZADD", "ranking:test", "1", "player1")).Return(mock.Result(mock.RedisInt64(1)))
	type fields struct {
		client rueidis.Client
		name   string
	}
	type args struct {
		ctx context.Context
		p   domain.Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test AddPlayer",
			fields: fields{
				client: client,
				name:   "test",
			},
			args: args{
				ctx: context.Background(),
				p: domain.Player{
					ID:      "player1",
					Ranking: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &RedisQueue{
				client: tt.fields.client,
				name:   tt.fields.name,
			}
			if err := q.AddPlayer(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("RedisQueue.AddPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisQueue_Make(t *testing.T) {
	p := domain.Player{
		ID:      "player1",
		Ranking: 79,
	}

	// player with ranking n
	nrBrackets := 10
	maxRanking := 100

	bracketInterval := maxRanking / nrBrackets
	slot := p.Ranking / bracketInterval
	fmt.Println("slot", slot, "bracketInterval", bracketInterval)

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		t.Fatal(err)
	}

	key := fmt.Sprintf("test:list:%d", slot)
	cmdPush := client.B().Rpush().Key(key).Element(testutil.NewID()).Build()
	cmdPush2 := client.B().Rpush().Key(key).Element(testutil.NewID()).Build()
	cmdPop := client.B().Lpop().Key(key).Count(2).Build()
	cmdPop2 := client.B().Lpop().Key(key).Count(2).Build()

	ctx, cancelMatchRequest := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer cancelMatchRequest()
	err = client.Do(ctx, cmdPush).Error()
	err = client.Do(context.Background(), cmdPush2).Error()

	if err != nil {
		t.Fatal(err)
	}

	r := client.DoMulti(context.Background(), cmdPop, cmdPop2)
	for _, v := range r {
		_, _ = v.AsStrSlice()
		// fmt.Println("v:", e, a)
	}

	// player with ranking n+1

}
