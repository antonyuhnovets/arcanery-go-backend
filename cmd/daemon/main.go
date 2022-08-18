package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	var l = &websocket.Lobby{
		Broadcast:  make(chan string),
		Register:   make(chan *websocket.Client),
		Unregister: make(chan *websocket.Client),
		Clients:    make(map[*websocket.Client]bool),
		Content:    "",
	}

	go l.Run()

	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", nil)
	// })

	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})
	// r.GET("/ws", func(c *gin.Context) {
	// 	websocket.ReadMsgsWS(c.Writer, c.Request)
	// })

	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(l, c.Writer, c.Request)
	})

	r.Run("localhost:8080")
}
