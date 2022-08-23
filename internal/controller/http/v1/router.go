package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	h := router.Group("v1/room")

	h.GET(":roomId", LoadFrontend)

	h.GET("ws/:roomId", StartWebsocket)
}
