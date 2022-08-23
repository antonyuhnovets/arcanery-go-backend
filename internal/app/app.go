package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hetonei/arcanery-go-backend/config"
	_ "github.com/hetonei/arcanery-go-backend/docs"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
	// "github.com/hetonei/arcanery-go-backend/pkg/httpserver"
)

// @title         Arcanery
// @version       0.0.1
// @description   Give me request and I'll give you a power
// @license.name  MIT
func Run(cfg *config.Config) {
	go websocket.H.Run()

	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.LoadHTMLFiles("index.html")

	v1.NewRouter(router)

	router.Run("0.0.0.0:8080")
}
