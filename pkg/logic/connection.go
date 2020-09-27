package logic

import (
	"net"
	"net/http"

	"github.com/egoon/hanabi-server/pkg/io"

	"github.com/egoon/hanabi-server/pkg/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func HandleConnection(conn net.Conn, games map[model.GameID]*model.Game, gameChan chan *model.Game) {
	defer conn.Close()
	ar := io.NewActionReader(conn)
	writer := io.NewJsonWriter(conn)
	var game *model.Game
	var playerID model.PlayerID
	for {
		action, err := ar.ReadAction()
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				log.Info("Connection timed out:", err)
				_, _ = writer.Write(model.Error{Err: http.StatusGatewayTimeout})
				break
			}
			log.Warn("Failed to read from client: ", err)
			_, _ = writer.Write(model.Error{Err: http.StatusBadRequest})
			break
		}
		if playerID == "" {
			if action.ActivePlayer == "" {
				playerID = model.PlayerID(uuid.New().String())
			} else {
				playerID = action.ActivePlayer
			}
		}
		if game == nil {
			err = ValidateAndCleanAction(action, nil)
			if err != nil {
				log.Info("validate action failed: ", err)

				_, err = writer.Write(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
				if err != nil {
					log.Warn("failed to send message to client: ", err)
				}
			} else {
				game, err = ConnectToGame(action, conn, games, gameChan)
				if err != nil {
					_, err = writer.Write(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
					if err != nil {
						log.Warn("failed to send message to client: ", err)
					}
				} else {
					_, err = writer.Write(model.GameState{Id: game.Id})
					if err != nil {
						log.Warn("failed to send message to client: ", err)
					}
				}
			}
		} else {
			err = ValidateAndCleanAction(action, game.State)
			if err != nil {
				log.Info("validate action failed: ", err)
				_, err = writer.Write(model.Error{Err: http.StatusBadRequest, Message: err.Error()})
				if err != nil {
					log.Warn("failed to send message to client: ", err)
				}
			} else {
				game.Actions <- action
			}
		}
	}
}
