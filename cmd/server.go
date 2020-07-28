package main

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/egoon/hanabi-server/pkg/logic"
	"github.com/egoon/hanabi-server/pkg/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	ln, err := net.Listen("tcp", ":579")
	if err != nil {
		log.Error("Failed to start server: ", err, ". Exiting\n")
	}
	games := map[model.GameID]*model.Game{}
	gameChan := make(chan *model.Game, 5)

	go logic.HandleNewGames(games, gameChan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warn("waiting for connection failed: ", err)
		} else {
			go handleConnection(conn, games, gameChan)
		}
	}
}

func handleConnection(conn net.Conn, games map[model.GameID]*model.Game, gameChan chan *model.Game) {
	buffer := make([]byte, 200)
	defer conn.Close()
	var game *model.Game
	var playerID model.PlayerID
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			log.Warn("set read deadline failed, game:", game.Id, ", player: ", playerID)
		}
		_, err = conn.Read(buffer)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				log.Info("Connection timed out:", err)
				msg, _ := json.Marshal(model.Error{Err: http.StatusGatewayTimeout})
				conn.Write(msg)
				break
			}
			log.Warn("Failed to read from client: ", err)
			msg, _ := json.Marshal(model.Error{Err: http.StatusBadRequest})
			_, err = conn.Write(msg)
			if err != nil {
				log.Warn("failed to send message to client: ", err)
			}
		}
		action := model.Action{}
		err = json.Unmarshal(buffer, &action)
		if err != nil {
			log.Info("failed to parse incoming action: ", err)
			msg, _ := json.Marshal(model.Error{Err: http.StatusBadRequest})
			_, err = conn.Write(msg)
			if err != nil {
				log.Warn("failed to send message to client: ", err)
			}
		}
		if playerID == "" {
			if action.ActivePlayer == "" {
				playerID = model.PlayerID(uuid.New().String())
			} else {
				playerID = action.ActivePlayer
			}
		}
		if game == nil {
			err = logic.ValidateAndCleanAction(&action, nil)
			if err != nil {
				log.Info("validate action failed: ", err)
				msg, _ := json.Marshal(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
				_, err = conn.Write(msg)
				if err != nil {
					log.Warn("failed to send message to client: ", err)
				}
			} else {
				game, err = logic.ConnectToGame(action, conn, games, gameChan)
				if err != nil {
					msg, _ := json.Marshal(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
					_, err = conn.Write(msg)
					if err != nil {
						log.Warn("failed to send message to client: ", err)
					}
				} else {
					msg, _ := json.Marshal(model.GameState{Id: game.Id})
					_, err = conn.Write(msg)
					if err != nil {
						log.Warn("failed to send message to client: ", err)
					}
				}
			}
		} else {
			err = logic.ValidateAndCleanAction(&action, game.State)
			if err != nil {
				log.Info("validate action failed: ", err)
				msg, _ := json.Marshal(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
				_, err = conn.Write(msg)
				if err != nil {
					log.Warn("failed to send message to client: ", err)
				}
			} else {
				game.Actions <- action
			}
		}
	}
}
