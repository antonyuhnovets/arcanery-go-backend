package v1

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

// @Summary     Load frontend with websocket
// @Description Load frontend and start chat room
// @ID          RoomFront
// @Tags  	    Frontend
// @Accept      json
// @Produce     html
// @Router      /{roomId} [get]
func LoadFrontend(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// @Summary     Start room with websocket
// @Description Load frontend and start chat room
// @ID          RoomBack
// @Tags  	    Backend
// @Accept      json
// @Produce     json
// @Router      /{roomId} [get]
func StartWebsocket(c *gin.Context) {
	id := c.Param("roomId")
	websocket.ServeWs(c.Writer, c.Request, id)
	log.Println(id)
}
