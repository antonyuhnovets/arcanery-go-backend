package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	h := router.Group("/v1")
	{
		newRoomRoutes(h)
	}
}
