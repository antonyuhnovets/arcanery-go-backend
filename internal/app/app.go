// Set up application before run

package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hetonei/arcanery-go-backend/config"
	_ "github.com/hetonei/arcanery-go-backend/docs"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	"github.com/hetonei/arcanery-go-backend/internal/middleware"
	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
)

// @title         Arcanery
// @version       0.0.1
// @description   Give me request and I'll give you a power
// @license.name  MIT
func Run(cfg *config.Config) {
	router := gin.New()

	// run hub
	go lobby.H.Run()

	// swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set routes
	v1.NewRouter(router)

	// use middlewares
	router.Use(middleware.Cors())

	// run on given host, port
	router.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
