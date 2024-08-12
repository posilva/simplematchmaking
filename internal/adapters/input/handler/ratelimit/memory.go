package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	rl "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

type KeyFunc func(*gin.Context) string

// Memory is an in memory rate limiter
type Memory struct {
	mw gin.HandlerFunc
}

// NewInMemory creates a new in memory rate limiter
func NewInMemory(rate time.Duration, limit uint) *Memory {
	store := rl.InMemoryStore(&rl.InMemoryOptions{
		Rate:  rate,
		Limit: limit,
	})
	mw := rl.RateLimiter(store, &rl.Options{
		ErrorHandler: defaultErrorHandler,
		KeyFunc:      defaultKeyFunc,
	})

	return &Memory{
		mw: mw,
	}
}

func NewInMemoryWithKeyFunc(rate time.Duration, limit uint, kf KeyFunc) *Memory {
	store := rl.InMemoryStore(&rl.InMemoryOptions{
		Rate:  rate,
		Limit: limit,
	})
	mw := rl.RateLimiter(store, &rl.Options{
		ErrorHandler: defaultErrorHandler,
		KeyFunc:      kf,
	})

	return &Memory{
		mw: mw,
	}
}

func defaultErrorHandler(ctx *gin.Context, info rl.Info) {
	ctx.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", info.Limit))
	ctx.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%v", info.RemainingHits))
	ctx.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", info.ResetTime.Unix()))
	ctx.AbortWithStatus(http.StatusTooManyRequests)
}

func defaultKeyFunc(ctx *gin.Context) string {
	fmt.Println(ctx.ClientIP())
	return "global"
}

// Handler returns the gin middleware
func (m *Memory) Handler() gin.HandlerFunc {
	return m.mw
}
