package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/posilva/simplematchmaking/internal/core/services"
)

// HTTPHandler is the HTTP Handler
type HTTPHandler struct {
	service ports.MatchmakingService
}

// NewHTTPHandler creates a new HTTP Handler
func NewHTTPHandler(srv ports.MatchmakingService) *HTTPHandler {
	return &HTTPHandler{
		service: srv,
	}
}

// HandleFindMatch handles the find match request
func (h *HTTPHandler) HandleFindMatch(ctx *gin.Context) {
	var in FindMatchInput
	err := ctx.Bind(&in)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid input")
		return
	}

	t, err := h.service.FindMatch(ctx.Request.Context(), "global", domain.Player{
		ID:      in.PlayerID,
		Ranking: in.Score,
	})

	if err != nil {
		// TODO: log error
		ctx.String(http.StatusInternalServerError, "failed to find a match")
		return
	}

	ctx.JSON(http.StatusOK, FindMatchOutput{
		TicketID: t.ID,
	})
}

// HandleCheckMatch handles the get match request
func (h *HTTPHandler) HandleCheckMatch(ctx *gin.Context) {
	ticketID := ctx.Params.ByName("ticketId")

	m, err := h.service.CheckMatch(ctx.Request.Context(), ticketID)
	if err != nil {
		if errors.Is(err, services.ErrMatchNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to get match with ticket: %v", ticketID))
		return
	}
	ctx.JSON(http.StatusOK, GetMatchOutput{
		MatchID: m.ID,
	})
}

// HandleCancelMatch handles the cancel match request
func (h *HTTPHandler) HandleCancelMatch(ctx *gin.Context) {
	ticketID := ctx.Params.ByName("ticketId")
	err := h.service.CancelMatch(ctx.Request.Context(), ticketID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to cancel match with ticket: %v", ticketID))
		return
	}
	ctx.Status(http.StatusNoContent)
}

// HandleRoot handles the root request
func (h *HTTPHandler) HandleRoot(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Simple Matchmaking Service")
}
