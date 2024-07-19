package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/config"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler"
	"github.com/posilva/simplematchmaking/internal/adapters/output/logging"
	"github.com/posilva/simplematchmaking/internal/adapters/output/matchmaker"
	"github.com/posilva/simplematchmaking/internal/adapters/output/repository"
	"github.com/redis/rueidis"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/services"
)

func Run() {
	r := gin.Default()

	service, err := createService()
	if err != nil {
		panic(fmt.Errorf("failed to create service instance: %v", err))
	}

	httpHandler := handler.NewHTTPHandler(service)
	r.GET("/", httpHandler.HandleRoot)
	api := r.Group("api/v1")

	api.PUT("/queue", httpHandler.HandleFindMatch)
	api.GET("/queue/:ticketId", httpHandler.HandleCheckMatch)
	api.DELETE("/queue/:ticketId", httpHandler.HandleCancelMatch)

	err = r.Run(config.GetAddr())
	if err != nil {
		panic(fmt.Errorf("failed to start the server %v", err))
	}

}

func createService() (*services.MatchmakingService, error) {
	logger := logging.NewSimpleLogger()

	rc, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{config.GetRedisAddr()}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %v", err)
	}

	mmCfg := domain.MatchmakerConfig{
		MaxPlayers: 2,
	}

	mm, err := matchmaker.NewRedisMatchmaker(rc, mmCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create matchmaker: %v", err)
	}

	repo := repository.NewRedisRepository(rc)

	return services.NewMatchmakingService(logger, repo, mm), nil
}
