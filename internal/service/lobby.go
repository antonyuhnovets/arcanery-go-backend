package service

// import(
// 	"github.com/hetonei/arcanery-go-backend/internal/domain/models"
// )

type Lobby struct {
	room RoomService
	game GameService
}

type RoomService interface {
	CreateRoom(string)
	DeleteRoom(string)
	ConnectToRoom(string)
	Disconnect(string)
}

type PlayerService interface {
	ChoseCards(...string)
	PlayCard(string)
	GetReady()
	ChangeCard(string, string)
}

type GameService interface {
	StartGame()
	EndTurn()
	EndRound()
	EndGame()
}

func CreateLobby(r Lobby) {
	return
}
