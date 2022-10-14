package http

import (
	"github.com/gin-gonic/gin"
)

// Router for room
func NewRouter(router *gin.Engine) {
	h := router.Group("v1/room")

	h.GET("/new", CreateRoom)

	h.GET("/:roomId", ConnectById)

	h.GET("ws/:roomId", ConnectWS)

	h.DELETE("/:roomId", DeleteRoomById)
}
