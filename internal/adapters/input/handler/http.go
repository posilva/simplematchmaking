package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	ctx.String(http.StatusOK, "OK")
}
