package health

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	hconfig "github.com/tavsec/gin-healthcheck/config"
)

func Setup(r *gin.Engine, redisClient rueidis.Client) {
	redisCheck := NewRedisCheck(redisClient)

	healthcheck.New(r, hconfig.DefaultConfig(), []checks.Check{
		redisCheck,
	})
}
