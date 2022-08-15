package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		websocket.ReadMsgsWS(c.Writer, c.Request, c.Request.Header)
	})

	r.Run("localhost:8000")
}
