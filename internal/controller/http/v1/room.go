package v1

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

func newRoomRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/room")
	l := websocket.NewLobby()
	go l.Run()
	{
		h.GET(":roomId", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
		h.GET("ws/:roomId", func(c *gin.Context) {
			id := c.Param("roomId")
			websocket.ServeWs(l, c.Writer, c.Request, id)
			log.Println(id)
		})
	}
}
