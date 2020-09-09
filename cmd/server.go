package main

import (
	"net"

	"github.com/egoon/hanabi-server/pkg/logic"
	"github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	ln, err := net.Listen("tcp", ":579")
	if err != nil {
		log.Error("Failed to start server: ", err, ". Exiting\n")
	}
	games := map[model.GameID]*model.Game{}
	gameChan := make(chan *model.Game, 5)

	// only HandleNewGames is allowed to add entries to the 'games' map
	go logic.HandleNewGames(games, gameChan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warn("waiting for connection failed: ", err)
		} else {
			go logic.HandleConnection(conn, games, gameChan)
		}
	}
}
