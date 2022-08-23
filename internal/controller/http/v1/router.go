package v1

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

func NewRouter(router *gin.Engine) {
	// h := router.Group("/v1")
	router.GET(":roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("ws/:roomId", func(c *gin.Context) {
		id := c.Param("roomId")
		websocket.ServeWs(c.Writer, c.Request, id)
		log.Println(id)
	})
}
