package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports/mocks"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMatchmakingService_FindMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	mmMock.EXPECT().AddPlayer(ctx, gomock.Any()).Return(nil)
	repoMock.EXPECT().ReservePlayerSlot(
		ctx, "player1", "queue1", gomock.Any()).Return(true, nil)
	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	ticket, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: "player1",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, ticket.ID)
}

func TestMatchmakingService_FindMatch_Exist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	repoMock.EXPECT().ReservePlayerSlot(
		ctx, "player1", "queue1", gomock.Any()).Return(false, nil)
	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	_, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: "player1",
	})

	assert.ErrorContains(t, err, "player already in the queue")
}

func TestMatchmakingService_FindMatch_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	repoMock.EXPECT().ReservePlayerSlot(
		ctx, "player1", "queue1", gomock.Any()).Return(false, fmt.Errorf("any error"))
	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	_, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: "player1",
	})

	assert.ErrorContains(t, err, "any error")
}
