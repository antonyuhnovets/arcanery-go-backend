package http

import (
	"github.com/gin-gonic/gin"

	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
	"github.com/hetonei/arcanery-go-backend/pkg/uuid"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

// @Summary     Create room
// @Description Start chat room
// @ID          RoomBack
// @Tags  	    Backend
// @Accept      json
// @Produce     json
// @Router      /new [post]
func CreateRoom(rh *lobby.EventHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.GenerateId()

		req := &lobby.Request{
			Event:  "create",
			Params: map[string]string{"room": id},
		}

		rh.MakeRequest(req)

		c.JSON(200, id)

	}
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

func ConnectWS(rh *lobby.EventHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("roomId")
		connId := uuid.GenerateId()

		err := rh.RoomCheck(id)

		if err != nil {
			c.JSON(404, err)
			return
		}

		ws, err := websocket.UpgradeToWs(c.Writer, c.Request)
		if err != nil {
			c.JSON(400, "Can't upgrade connection")
			return
		}

		conn := websocket.MakeConnection(ws)

		req := lobby.Request{
			Event: "open",
			Params: map[string]string{
				"room":     id,
				"connId":   connId,
				"connType": "websocket",
			},
			Data: map[string]interface{}{
				"conn": conn,
			},
		}

		rh.MakeRequest(&req)
	}
}

// @Summary     Delete Room
// @Description Remove room by id
// @ID          rmRoom
// @Tags  	    del
// @Accept      json
// @Produce     html
// @Router      /{roomId} [delete]
func DeleteRoomById(c *gin.Context) {
	// id := c.Param("roomId")
	// srv, ws := GetServices(c, id, M)

	// srv.DeleteRoom(id)
}
