package model

import "net"

type Game struct {
	Id          GameID
	Connections map[PlayerID]net.Conn
	Actions     chan *Action
	State       *GameState
}
