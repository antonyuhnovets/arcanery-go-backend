package app

import (
	"fmt"
	// "net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hetonei/arcanery-go-backend/config"
	_ "github.com/hetonei/arcanery-go-backend/docs"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	"github.com/hetonei/arcanery-go-backend/internal/middleware"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

// @title         Arcanery
// @version       0.0.1
// @description   Give me request and I'll give you a power
// @license.name  MIT
func Run(cfg *config.Config) {
	router := gin.New()

	go websocket.H.Run()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.LoadHTMLFiles("index.html")

	v1.NewRouter(router)

	router.Use(middleware.Cors())

	router.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
