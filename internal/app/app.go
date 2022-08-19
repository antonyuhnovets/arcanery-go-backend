package app

import (
	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/config"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	// "github.com/hetonei/arcanery-go-backend/pkg/httpserver"
)

func Run(cfg *config.Config) {
	handler := gin.New()
	handler.LoadHTMLFiles("index.html")
	v1.NewRouter(handler)

	handler.Run("localhost:8000")
}
