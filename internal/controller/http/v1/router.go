package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	h := router.Group("v1/room")

	h.GET("/new", CreateRoom)

	h.GET("/rm/:roomId", DeleteRoomById)

	h.GET(":roomId", ConnectById)

	h.GET("ws/:roomId", ConnectWS)
}
