package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetUpRouter(mw gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(mw)
	return router
}

func request() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	return w, req
}

func TestRatelimitInMemory(t *testing.T) {
	r := SetUpRouter(NewInMemory(time.Second, 2).Handler())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w, req := request()
	r.ServeHTTP(w, req)
	w, req = request()
	r.ServeHTTP(w, req)
	w, req = request()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRatelimitInMemoryWithKeyFunc(t *testing.T) {
	r := SetUpRouter(NewInMemoryWithKeyFunc(time.Second, 2, func(c *gin.Context) string {
		return c.Request.Host
	}).Handler())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w, req := request()
	r.ServeHTTP(w, req)
	w, req = request()
	r.ServeHTTP(w, req)
	w, req = request()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusTooManyRequests, w.Code)
}
