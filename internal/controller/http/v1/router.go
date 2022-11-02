package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
)

// Router for room
func NewRouter(router *gin.Engine, rh *lobby.EventHub) {
	h := router.Group("v1/room")

	h.GET("/new", CreateRoom(rh))

	h.GET("/:roomId", ConnectById)

	h.GET("ws/:roomId", ConnectWS(rh))

	h.DELETE("/:roomId", DeleteRoomById)
}
