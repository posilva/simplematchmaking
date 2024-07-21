package queues

import (
	"context"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
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
