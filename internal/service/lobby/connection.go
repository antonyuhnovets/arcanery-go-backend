package lobby

import (
	"log"
)

type Connection interface {
	SetConnectionId(string)
	GetConnectionId() string
	Listen(func([]byte))
	SendMsg(interface{})
	Close()
}

type ConnectionManager struct {
	connType string
	params   map[string]interface{}
}

func NewConnectionManager(connType string, data map[string]interface{}) Connection {
	log.Printf("%s connection type", connType)
	return ConnectionManager{connType: connType, params: data}.GetConnection()
}

func (cm ConnectionManager) GetConnection() Connection {
	return cm.params["conn"].(Connection)
}
