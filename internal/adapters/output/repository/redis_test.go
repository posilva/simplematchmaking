package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/ports/mocks"
	"github.com/redis/rueidis/mock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReservePlayerSlotOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	slot := "slot1"
	playerID := "player1"

	clt := mock.NewClient(ctrl)
	ticketID := ksuid.New().String()
	key := fmt.Sprintf("playerslot:%s:%s", slot, playerID)
	value := fmt.Sprintf("status:reserved:ticket:%s", ticketID)
	clt.EXPECT().Do(gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
	).Return(
		mock.Result(mock.RedisInt64(1)),
	)

	codec := mocks.NewMockCodec(ctrl)
	repo := NewRedisRepository(clt, codec)

	ok, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	assert.NoError(t, err)
	assert.True(t, ok)

}

func TestReservePlayerSlotExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	slot := "slot1"
	playerID := "player1"
	clt := mock.NewClient(ctrl)
	ticketID := ksuid.New().String()
	key := fmt.Sprintf("playerslot:%s:%s", slot, playerID)
	value := fmt.Sprintf("status:reserved:ticket:%s", ticketID)

	clt.EXPECT().Do(gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
	).Return(
		mock.Result(mock.RedisInt64(0)),
	)

	codec := mocks.NewMockCodec(ctrl)
	repo := NewRedisRepository(clt, codec)

	ok, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	assert.NoError(t, err)
	assert.False(t, ok)

}

func TestReservePlayerSlotError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	errMsg := "unexpected test error"

	ctx := context.Background()
	slot := "slot1"
	playerID := "player1"
	clt := mock.NewClient(ctrl)
	ticketID := ksuid.New().String()
	key := fmt.Sprintf("playerslot:%s:%s", slot, playerID)
	value := fmt.Sprintf("status:reserved:ticket:%s", ticketID)

	clt.EXPECT().Do(gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
	).Return(
		mock.ErrorResult(fmt.Errorf(errMsg)),
	)
	codec := mocks.NewMockCodec(ctrl)
	repo := NewRedisRepository(clt, codec)

	ok, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	assert.ErrorContains(t, err, errMsg)
	assert.False(t, ok)

}
