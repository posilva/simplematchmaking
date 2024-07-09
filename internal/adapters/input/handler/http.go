package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
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

// Handle handles the GET / endpoint
func (h *HTTPHandler) Handle(ctx *gin.Context) {
	t, err := h.service.FindMatch(domain.Player{
		ID: "1",
	})
	if err != nil {
		// TODO: log error
		ctx.String(http.StatusInternalServerError, "failed to find a match")
		return
	}

	ctx.JSON(http.StatusOK, t)
}
