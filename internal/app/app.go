package app

import (
	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/config"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
	// "github.com/hetonei/arcanery-go-backend/pkg/httpserver"
)

func Run(cfg *config.Config) {
	go websocket.H.Run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")
	v1.NewRouter(router)

	router.Run("localhost:8000")
}
