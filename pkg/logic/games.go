package logic

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/egoon/hanabi-server/pkg/model"
)

func HandleNewGames(games map[model.GameID]*model.Game, gameChan chan *model.Game) {
	for {
		game := <-gameChan
		if games[game.Id] != nil {
			for _, conn := range game.Connections {
				conn.Close()
			}
		} else {
			games[game.Id] = game
		}
	}
}

func ConnectToGame(action model.Action, conn net.Conn, games map[model.GameID]*model.Game, gameChan chan *model.Game) (*model.Game, error) {
	playerID := action.ActivePlayer
	switch action.Type {
	case "create":
		_, ok := games[action.GameID]
		if ok {
			return nil, fmt.Errorf("cannot create game. game already exists")
		}
		actions := make(chan model.Action, 5)
		game := &model.Game{
			Id: action.GameID,
			Connections: map[model.PlayerID]net.Conn{
				playerID: conn,
			},
			Actions: actions,
		}
		// send game to async game handler function HandleNewGames
		gameChan <- game
		// create an async func to handle the new games actions
		go func() {
			HandleGameActions(game, model.CreateDeck())
			delete(games, game.Id)
		}()
		game.Actions <- model.Action{
			Type:         "join",
			ActivePlayer: playerID,
		}
		return game, nil
	case "join":
		game := games[action.GameID]
		if game == nil {
			msg, _ := json.Marshal(model.Error{Err: http.StatusNotFound})
			conn.Write(msg)
			return nil, fmt.Errorf("cannot join game. game does not exist")
		}
		previousConn := game.Connections[playerID]
		if previousConn != nil {
			previousConn.Close()
		} else if len(game.Connections) > 4 {
			msg, _ := json.Marshal(model.Error{Err: http.StatusPreconditionFailed})
			conn.Write(msg)
			return nil, fmt.Errorf("cannot join game. too many connections")
		}
		game.Connections[playerID] = conn
		game.Actions <- action
		return game, nil
	default:
		msg, _ := json.Marshal(model.Error{Err: http.StatusConflict})
		conn.Write(msg)
		return nil, fmt.Errorf("invalid action: %s. not in a game", action.Type)
	}
}
