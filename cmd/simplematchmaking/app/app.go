package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/config"
	"github.com/posilva/simplematchmaking/internal/adapters/input/handler"

	"github.com/posilva/simplematchmaking/internal/core/services"
)

func Run() {
	r := gin.Default()

	service, err := createService()
	if err != nil {
		panic(fmt.Errorf("failed to create service instance: %v", err))
	}

	httpHandler := handler.NewHTTPHandler(service)
	r.GET("/", httpHandler.Handle)
	api := r.Group("api/v1")

	api.GET("/mm", httpHandler.Handle)

	err = r.Run(config.GetAddr())
	if err != nil {
		panic(fmt.Errorf("failed to start the server %v", err))
	}

}

func createService() (*services.MatchmakingService, error) {
	return services.NewMatchmakingService(), nil
}
