package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports/mocks"
	"github.com/posilva/simplematchmaking/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMatchmakingService_FindMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	pID := testutil.NewID()

	mmMock.EXPECT().Subscribe(gomock.Any()).Return()
	mmMock.EXPECT().AddPlayer(gomock.Any(), gomock.Any()).Return(nil)

	rps := repoMock.EXPECT().ReservePlayerSlot(ctx, pID, "queue1", gomock.Any())
	var expectedID string
	rps.Do(func(ctx context.Context, playerID string, slot string, ticketID string) (string, error) {
		expectedID = ticketID
		rps.Return(ticketID, nil)
		return ticketID, nil
	})

	repoMock.EXPECT().UpdateTicket(gomock.Any(), gomock.Any()).Return(nil)

	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	ticket, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: pID,
	})

	require.NoError(t, err)
	require.Equal(t, expectedID, ticket.ID)
}

func TestMatchmakingService_FindMatch_Exist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tID := testutil.NewID()

	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	mmMock.EXPECT().Subscribe(gomock.Any()).Return()
	// mmMock.EXPECT().AddPlayer(gomock.Any(), gomock.Any()).Return(nil)

	// repoMock.EXPECT().UpdateTicket(gomock.Any(), gomock.Any()).Return(nil)
	repoMock.EXPECT().ReservePlayerSlot(
		ctx, "player1", "queue1", gomock.Any()).Return(tID, nil)

	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	expectedID, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: "player1",
	})

	require.NoError(t, err)
	require.Equal(t, tID, expectedID.ID)
}

func TestMatchmakingService_FindMatch_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoMock := mocks.NewMockRepository(ctrl)
	mmMock := mocks.NewMockMatchmaker(ctrl)
	log := testutil.NewLogger(t)

	mmMock.EXPECT().Subscribe(gomock.Any()).Return()
	repoMock.EXPECT().ReservePlayerSlot(
		ctx, "player1", "queue1", gomock.Any()).Return("", fmt.Errorf("any error"))
	s := NewMatchmakingService(
		log,
		repoMock,
		mmMock,
	)
	_, err := s.FindMatch(ctx, "queue1", domain.Player{
		ID: "player1",
	})

	require.ErrorContains(t, err, "any error")
}
