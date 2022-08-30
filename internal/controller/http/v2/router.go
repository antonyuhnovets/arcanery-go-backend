package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func NewRouter(router *gin.Engine, socket *socketio.Server) {
	h := router.Group("v2")

	h.GET("/socket.io/*any", gin.WrapH(socket))
	h.POST("/socket.io/*any", gin.WrapH(socket))
	router.StaticFS("/public", http.Dir("./asset"))
}
