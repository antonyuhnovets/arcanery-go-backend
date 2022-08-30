package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	Server *socketio.Server
}

type Client struct {
	Conn *socketio.Conn
}

func NewServer() *Server {
	server := socketio.NewServer(nil)
	return &Server{server}
}

func (s *Server) Run() {
	if err := s.Server.Serve(); err != nil {
		log.Fatalf("socketio listen error: %s\n", err)
	}
}

func (s *Server) Stop() {
	s.Server.Close()
}

func (s *Server) ConnectClient(room string) {
	s.Server.OnConnect("/", func(c socketio.Conn) error {
		c.Join(room)
		c.SetContext(room)
		log.Println(c.ID(), " user connected to :", room)
		return nil
	})
}

func (s *Server) ListenMsgs() {
	s.Server.OnEvent("/", "notice", func(c socketio.Conn, msg string) {
		log.Println("notice:", msg)
		c.Emit("reply", "have "+msg)
		s.Server.BroadcastToRoom("/chat", c.Context().(string), "msg")
	})

	s.Server.OnEvent("/chat", "msg", func(c socketio.Conn, msg string) string {
		// c.SetContext(msg)
		s.Server.BroadcastToRoom("/chat", c.Context().(string), "msg")
		return "recv " + msg
	})

	s.Server.OnEvent("/", "bye", func(c socketio.Conn) string {
		last := c.Context().(string)
		c.Emit("bye", last)
		c.Close()
		return last
	})

	s.Server.OnError("/", func(c socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

}

func (s *Server) DisconnectClient() {
	s.Server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		log.Println("closed", reason)
	})
}
