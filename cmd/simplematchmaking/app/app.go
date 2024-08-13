package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/config"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler/health"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler/shutdown"
	"github.com/posilva/simplematchmaking/internal/adapters/output/lock"
	"github.com/posilva/simplematchmaking/internal/adapters/output/logging"
	"github.com/posilva/simplematchmaking/internal/adapters/output/queues"
	"github.com/posilva/simplematchmaking/internal/adapters/output/repository"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	configMM "github.com/posilva/simplematchmaking/internal/core/services/config"
	"github.com/redis/rueidis"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/domain/codecs"
	"github.com/posilva/simplematchmaking/internal/core/services"
)

// Run starts the application
func Run() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	rc, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{config.GetRedisAddr()},
	})
	if err != nil {
		panic(fmt.Errorf("failed to create redis client: %v", err))
	}

	service, err := createService(rc)
	if err != nil {
		panic(fmt.Errorf("failed to create service instance: %v", err))
	}

	health.Setup(r, rc)

	httpHandler := handler.NewHTTPHandler(service)

	r.GET("/", httpHandler.HandleRoot)
	api := r.Group("api/v1")

	api.PUT("/queue/:queue", httpHandler.HandleFindMatch)
	api.GET("/ticket/:ticketId", httpHandler.HandleCheckMatch)
	api.DELETE("/ticket/:ticketId", httpHandler.HandleCancelMatch)

	shut := shutdown.New()
	defer shut.Stop()

	shut.Start(
		func() {
			err = r.Run(config.GetAddr())
			if err != nil {
				panic(fmt.Errorf("failed to start the server %v", err))
			}
		})
}

func createService(rc rueidis.Client) (*services.MatchmakingService, error) {
	logger := logging.NewSimpleLogger()

	codec := codecs.NewJSONCodec()

	envVarConfig := configMM.NewEnvVar()
	err := envVarConfig.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	lock, err := lock.NewRedisLock(rc, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create matchmaker: %v", err)
	}

	repo := repository.NewRedisRepository(rc, codec, logger)

	matchmakers := make(map[string]ports.Matchmaker)

	for mmName, mmCfg := range envVarConfig.Get().Matchmakers {
		if qCfg, ok := envVarConfig.Get().Queues[mmCfg.Name]; ok {
			q := queues.NewRedisQueue(rc, qCfg, codec, lock)
			mm, err := services.NewMatchmaker(q, mmCfg, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to create matchmaker with name '%v': %v", q.Name(), err)
			}
			matchmakers[mmName] = mm
			continue
		}
		return nil, fmt.Errorf("queue with name '%v' not found", mmCfg.Name)
	}

	if len(matchmakers) == 0 {
		mm, name, err := defaultMatchmaker(rc, logger, codec, lock)
		if err != nil {
			return nil, fmt.Errorf("failed to create default matchmaker: %v", err)
		}
		matchmakers[name] = mm
	}

	return services.NewMatchmakingService(logger, repo, matchmakers), nil
}

func defaultMatchmaker(rc rueidis.Client, logger ports.Logger, codec ports.Codec, lock ports.Lock) (ports.Matchmaker, string, error) {
	name := "default"
	mmCfg := domain.MatchmakerConfig{
		Name:            name,
		IntervalSecs:    5,
		MakeTimeoutSecs: 4,
	}

	qConfig := domain.QueueConfig{
		MaxPlayers:     2,
		NrBrackets:     100,
		MinRanking:     1,
		MaxRanking:     1000,
		MakeIterations: 3,
		Name:           name,
	}

	queue := queues.NewRedisQueue(rc, qConfig, codec, lock)
	mm, err := services.NewMatchmaker(queue, mmCfg, logger)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create matchmaker: %v", err)
	}
	return mm, name, nil
}
