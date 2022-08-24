package v1

import (
	"github.com/gin-gonic/gin"

	rs "github.com/hetonei/arcanery-go-backend/internal/service/room_service"
	"github.com/hetonei/arcanery-go-backend/pkg/uuid"
)

// @Summary     Create room
// @Description Start chat room
// @ID          RoomBack
// @Tags  	    Backend
// @Accept      json
// @Produce     json
// @Router      /new [get]
func CreateRoom(c *gin.Context) {
	srv := rs.RegisterClient(c.Writer, c.Request)
	id := uuid.GenerateId()

	srv.CreateRoom(id)

}

// @Summary     Load frontend with websocket
// @Description Load frontend and start chat room
// @ID          RoomFront
// @Tags  	    Frontend
// @Accept      json
// @Produce     html
// @Router      /{roomId} [get]
func ConnectById(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func ConnectWS(c *gin.Context) {
	id := c.Param("roomId")
	srv := rs.RegisterClient(c.Writer, c.Request)

	srv.ConnectToRoom(id)
}

// @Summary     Delete Room
// @Description Remove room by id
// @ID          rmRoom
// @Tags  	    del
// @Accept      json
// @Produce     html
// @Router      /rm/{roomId} [get]
func DeleteRoomById(c *gin.Context) {
	id := c.Param("roomId")
	srv := rs.RegisterClient(c.Writer, c.Request)

	srv.DeleteRoom(id)
}
