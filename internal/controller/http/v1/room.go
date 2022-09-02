package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/internal/service"
	"github.com/hetonei/arcanery-go-backend/pkg/uuid"
)

// @Summary     Create room
// @Description Start chat room
// @ID          RoomBack
// @Tags  	    Backend
// @Accept      json
// @Produce     json
// @Router      /new [post]
func CreateRoom(c *gin.Context) {
	srv := service.RegisterRequest(c.Writer, c.Request)
	id := uuid.GenerateId()

	srv.CreateRoom(id)
	c.JSON(200, id)
}

// @Summary     Load frontend with websocket
// @Description Load frontend and start chat room
// @ID          RoomFront
// @Tags  	    Frontend
// @Accept      json
// @Produce     html
// @Router      /{roomId} [get]
func ConnectById(c *gin.Context) {
}

func ConnectWS(c *gin.Context) {
	id := c.Param("roomId")
	srv := service.RegisterRequest(c.Writer, c.Request)

	srv.ConnectToRoom(id)
}

// @Summary     Delete Room
// @Description Remove room by id
// @ID          rmRoom
// @Tags  	    del
// @Accept      json
// @Produce     html
// @Router      /{roomId} [delete]
func DeleteRoomById(c *gin.Context) {
	id := c.Param("roomId")
	srv := service.RegisterRequest(c.Writer, c.Request)

	srv.DeleteRoom(id)
}
