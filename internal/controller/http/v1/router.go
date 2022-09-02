package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	h := router.Group("v1/room")

	h.POST("/new", CreateRoom)

	h.GET("/:roomId", ConnectById)

	h.GET("ws/:roomId", ConnectWS)

	h.DELETE("/:roomId", DeleteRoomById)
}
