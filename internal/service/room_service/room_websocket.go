package websocket

import (
	"net/http"

	"github.com/hetonei/arcanery-go-backend/internal/service"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type ClientWS struct {
	id  string
	w   http.ResponseWriter
	req *http.Request
}

func GetClientService(w http.ResponseWriter, r *http.Request) service.ClientService {
	return ClientWS{
		id:  r.RemoteAddr,
		w:   w,
		req: r,
	}
}

func (cws ClientWS) ConnectToRoom(roomId string) {
	websocket.ServeWs(cws.w, cws.req, roomId)
}

func ConnectWSClient()
