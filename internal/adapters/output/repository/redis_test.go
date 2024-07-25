package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/domain/codecs"
	"github.com/posilva/simplematchmaking/internal/core/ports/mocks"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
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
	clt.EXPECT().DoMulti(
		gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
		mock.Match("GET", key),
	).Return(
		[]rueidis.RedisResult{
			mock.Result(mock.RedisNil()),
			mock.Result(mock.RedisString(value)),
		},
	)

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	tID, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	require.NoError(t, err)
	require.Equal(t, ticketID, tID)

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

	clt.EXPECT().DoMulti(
		gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
		mock.Match("GET", key),
	).Return(
		[]rueidis.RedisResult{
			mock.Result(mock.RedisNil()),
			mock.Result(mock.RedisString(value)),
		},
	)
	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	tID, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	require.NoError(t, err)
	require.Equal(t, ticketID, tID)

}

func TestReservePlayerSlotErrorSET(t *testing.T) {
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

	clt.EXPECT().DoMulti(
		gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
		mock.Match("GET", key),
	).Return(
		[]rueidis.RedisResult{
			mock.ErrorResult(fmt.Errorf(errMsg)),
			mock.ErrorResult(fmt.Errorf(errMsg)),
		},
	)
	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	tID, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	require.ErrorContains(t, err, ErrFailedToReservePlayerSlot.Error())
	require.Empty(t, tID)

}
func TestReservePlayerSlotErrorGET(t *testing.T) {
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

	clt.EXPECT().DoMulti(
		gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
		mock.Match("GET", key),
	).Return(
		[]rueidis.RedisResult{
			mock.Result(mock.RedisString("OK")),
			mock.ErrorResult(fmt.Errorf(errMsg)),
		},
	)
	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	tID, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	require.ErrorContains(t, err, ErrFailedToReservePlayerSlot.Error())
	require.Empty(t, tID)

}

func TestReservePlayerSlotErrorNeverHappens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	slot := "slot1"
	playerID := "player1"
	clt := mock.NewClient(ctrl)
	ticketID := ksuid.New().String()
	key := fmt.Sprintf("playerslot:%s:%s", slot, playerID)
	value := fmt.Sprintf("status:reserved:ticket:%s", ticketID)

	clt.EXPECT().DoMulti(
		gomock.Any(),
		mock.Match("SET", key, value, "NX", "EX", "60"),
		mock.Match("GET", key),
	).Return(
		[]rueidis.RedisResult{
			mock.Result(mock.RedisString("OK")),
			mock.Result(mock.RedisNil()),
		},
	)
	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	tID, err := repo.ReservePlayerSlot(ctx, playerID, slot, ticketID)
	require.ErrorContains(t, err, ErrSomethingOddHappened.Error())
	require.Empty(t, tID)

}
func TestDeletePlayerSlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)

	slot := "slot1"
	playerID := testutil.NewID()

	key := fmt.Sprintf("playerslot:%s:%s", slot, playerID)
	clt.EXPECT().Do(
		gomock.Any(),
		mock.Match("DEL", key),
	).Return(
		mock.Result(
			mock.RedisString("OK")),
	)

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewRedisRepository(clt, codec, logger)

	err := repo.DeletePlayerSlot(context.Background(), playerID, slot)
	require.NoError(t, err)

}

func TestGetTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GET", key)).
		Return(mock.Result(mock.RedisString("OK")))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Decode(gomock.Any(), gomock.Any()).Return(nil)

	repo := NewRedisRepository(clt, codec, logger)

	_, err := repo.GetTicket(context.Background(), tID)
	require.NoError(t, err)
}

func TestGetTicket_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GET", key)).
		Return(mock.Result(mock.RedisNil()))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	repo := NewRedisRepository(clt, codec, logger)

	_, err := repo.GetTicket(context.Background(), tID)
	require.ErrorContains(t, err, ErrTicketNotFound.Error())
}

func TestGetTicket_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	errMsg := "unexpected test error"
	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GET", key)).
		Return(mock.ErrorResult(fmt.Errorf(errMsg)))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	repo := NewRedisRepository(clt, codec, logger)

	_, err := repo.GetTicket(context.Background(), tID)
	require.ErrorContains(t, err, ErrFailedToGetTicket.Error())
}

func TestGetTicket_ErrorDecode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	errMsg := "unexpected test error"

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GET", key)).
		Return(mock.Result(mock.RedisString("OK")))

	codec := mocks.NewMockCodec(ctrl)
	codec.EXPECT().Decode(gomock.Any(), gomock.Any()).Return(fmt.Errorf(errMsg))
	logger := mocks.NewMockLogger(ctrl)

	repo := NewRedisRepository(clt, codec, logger)

	_, err := repo.GetTicket(context.Background(), tID)
	require.ErrorContains(t, err, ErrFailedToDecodeTicket.Error())
}

func TestDeleteTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	ticketSt := domain.TicketStatus{
		ID: tID,
	}

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GETDEL", key)).
		Return(mock.Result(mock.RedisString("OK")))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	decode := codec.EXPECT().Decode(gomock.Any(), gomock.Any())
	decode.Do(
		func(b []byte, v interface{}) error {
			*v.(*domain.TicketStatus) = ticketSt
			return nil
		}).Return(nil)

	repo := NewRedisRepository(clt, codec, logger)

	st, err := repo.DeleteTicket(context.Background(), tID)
	fmt.Println("Status:", st)
	require.NoError(t, err)
	require.Equal(t, st.ID, ticketSt.ID)
}

func TestDeleteTicket_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GETDEL", key)).
		Return(mock.Result(mock.RedisNil()))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	repo := NewRedisRepository(clt, codec, logger)
	_, err := repo.DeleteTicket(context.Background(), tID)

	require.ErrorContains(t, err, ErrTicketNotFound.Error())
}

func TestDeleteTicket_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	errMsg := "unexpected test error"

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GETDEL", key)).
		Return(mock.ErrorResult(fmt.Errorf(errMsg)))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	repo := NewRedisRepository(clt, codec, logger)

	_, err := repo.DeleteTicket(context.Background(), tID)

	require.ErrorContains(t, err, ErrFailedToDeleteTicket.Error())

}

func TestDeleteTicket_ErrorDecode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()

	key := fmt.Sprintf("ticket:%s", tID)
	clt.EXPECT().Do(gomock.Any(), mock.Match("GETDEL", key)).
		Return(mock.Result(mock.RedisString("OK")))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Decode(gomock.Any(), gomock.Any()).Return(fmt.Errorf("unexpected test error"))
	repo := NewRedisRepository(clt, codec, logger)
	_, err := repo.DeleteTicket(context.Background(), tID)

	require.ErrorContains(t, err, ErrFailedToDecodeTicket.Error())
}

func TestUpdateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	ticketSt := domain.TicketStatus{
		ID: tID,
	}

	key := fmt.Sprintf("ticket:%s", tID)
	jsonCodec := codecs.NewJSONCodec()

	enc, err := jsonCodec.Encode(ticketSt)
	value := string(enc)
	require.NoError(t, err)

	clt.EXPECT().Do(gomock.Any(), mock.Match("SET", key, value)).
		Return(mock.Result(mock.RedisString(value)))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Encode(gomock.Any()).Return(enc, nil)

	repo := NewRedisRepository(clt, codec, logger)

	err = repo.UpdateTicket(context.Background(), ticketSt)
	require.NoError(t, err)
}

func TestUpdateTicket_ErrorEncode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	ticketSt := domain.TicketStatus{
		ID: tID,
	}

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Encode(gomock.Any()).Return(nil, fmt.Errorf("unexpected test error"))

	repo := NewRedisRepository(clt, codec, logger)

	err := repo.UpdateTicket(context.Background(), ticketSt)
	require.ErrorContains(t, err, ErrFailedToEncodeTicket.Error())
}

func TestUpdateTicket_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	ticketSt := domain.TicketStatus{
		ID: tID,
	}

	key := fmt.Sprintf("ticket:%s", tID)
	jsonCodec := codecs.NewJSONCodec()

	enc, err := jsonCodec.Encode(ticketSt)
	value := string(enc)
	require.NoError(t, err)

	clt.EXPECT().Do(gomock.Any(), mock.Match("SET", key, value)).
		Return(mock.ErrorResult(fmt.Errorf("unexpected test error")))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Encode(gomock.Any()).Return(enc, nil)

	repo := NewRedisRepository(clt, codec, logger)

	err = repo.UpdateTicket(context.Background(), ticketSt)
	require.ErrorContains(t, err, ErrFailedToUpdateTicket.Error())
}

func TestUpdateTicket_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clt := mock.NewClient(ctrl)
	tID := testutil.NewID()
	ticketSt := domain.TicketStatus{
		ID: tID,
	}

	key := fmt.Sprintf("ticket:%s", tID)
	jsonCodec := codecs.NewJSONCodec()

	enc, err := jsonCodec.Encode(ticketSt)
	value := string(enc)
	require.NoError(t, err)

	clt.EXPECT().Do(gomock.Any(), mock.Match("SET", key, value)).
		Return(mock.Result(mock.RedisNil()))

	codec := mocks.NewMockCodec(ctrl)
	logger := mocks.NewMockLogger(ctrl)

	codec.EXPECT().Encode(gomock.Any()).Return(enc, nil)

	repo := NewRedisRepository(clt, codec, logger)

	err = repo.UpdateTicket(context.Background(), ticketSt)
	require.ErrorContains(t, err, ErrTicketNotFound.Error())
}
