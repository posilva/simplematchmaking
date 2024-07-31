package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/config"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler"
	"github.com/posilva/simplematchmaking/internal/adapters/output/lock"
	"github.com/posilva/simplematchmaking/internal/adapters/output/logging"
	"github.com/posilva/simplematchmaking/internal/adapters/output/queues"
	"github.com/posilva/simplematchmaking/internal/adapters/output/repository"
	"github.com/redis/rueidis"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/domain/codecs"
	"github.com/posilva/simplematchmaking/internal/core/services"
)

// Run starts the application
func Run() {
	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())

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
	lock, err := lock.NewRedisLock(rc, 1)

	codec := codecs.NewJSONCodec()

	mmCfg := domain.MatchmakerConfig{
		Name:            "main",
		IntervalSecs:    5,
		MakeTimeoutSecs: 4,
	}

	qConfig := domain.QueueConfig{
		MaxPlayers:     2,
		NrBrackets:     100,
		MinRanking:     1,
		MaxRanking:     1000,
		MakeIterations: 3,
		Name:           "global",
	}

	queue := queues.NewRedisQueue(rc, qConfig, codec, lock)
	mm, err := services.NewMatchmaker(queue, mmCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create matchmaker: %v", err)
	}

	repo := repository.NewRedisRepository(rc, codec, logger)

	return services.NewMatchmakingService(logger, repo, mm), nil
}
