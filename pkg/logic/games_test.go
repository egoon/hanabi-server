package logic

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egoon/hanabi-server/pkg/model"
)

func TestConnectToGame(t *testing.T) {
	testCases := []struct {
		description        string
		action             model.Action
		conn               *MockConn
		games              map[model.GameID]*model.Game
		gameChan           chan *model.Game
		expectedGame       *model.Game
		expectedErr        error
		expectedConnClosed bool
	}{
		{
			description: "Join empty game - ok",
			action: model.Action{
				Type:         model.ActionJoin,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn: &MockConn{BytesWritten: make(chan []byte, 5)},
			games: map[model.GameID]*model.Game{"ticTacToe": {
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{},
				Actions:     make(chan model.Action, 5),
			}},
			gameChan: make(chan *model.Game, 2),
			expectedGame: &model.Game{
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{"Top": &MockConn{}},
				Actions:     make(chan model.Action, 5),
			},
		},
		{
			description: "Join absent game - fail",
			action: model.Action{
				Type:         model.ActionJoin,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn:        &MockConn{BytesWritten: make(chan []byte, 5)},
			games:       map[model.GameID]*model.Game{},
			gameChan:    make(chan *model.Game, 2),
			expectedErr: fmt.Errorf("cannot join game. game does not exist"),
		},
		{
			description: "Re-Join game - ok",
			action: model.Action{
				Type:         model.ActionJoin,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn: &MockConn{BytesWritten: make(chan []byte, 5)},
			games: map[model.GameID]*model.Game{"ticTacToe": {
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{"Top": &MockConn{}},
				Actions:     make(chan model.Action, 5),
			}},
			gameChan: make(chan *model.Game, 2),
			expectedGame: &model.Game{
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{"Top": &MockConn{}},
				Actions:     make(chan model.Action, 5),
			},
		},
		{
			description: "Join full game - fail",
			action: model.Action{
				Type:         model.ActionJoin,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn: &MockConn{BytesWritten: make(chan []byte, 5)},
			games: map[model.GameID]*model.Game{"ticTacToe": {
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{"Bottom": nil, "Strange": nil, "Charm": nil, "Up": nil, "Down": nil},
				Actions:     make(chan model.Action, 5),
			}},
			gameChan:    make(chan *model.Game, 2),
			expectedErr: fmt.Errorf("cannot join game. too many connections"),
		},
		{
			description: "Create game - ok",
			action: model.Action{
				Type:         model.ActionCreate,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn:     &MockConn{BytesWritten: make(chan []byte, 5)},
			games:    map[model.GameID]*model.Game{},
			gameChan: make(chan *model.Game, 2),
			expectedGame: &model.Game{
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{"Top": &MockConn{}},
				Actions:     make(chan model.Action, 5),
			},
		},
		{
			description: "Create existing game - fail",
			action: model.Action{
				Type:         model.ActionCreate,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn: &MockConn{BytesWritten: make(chan []byte, 5)},
			games: map[model.GameID]*model.Game{"ticTacToe": {
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{},
				Actions:     make(chan model.Action, 5),
			}},
			gameChan:    make(chan *model.Game, 2),
			expectedErr: fmt.Errorf("cannot create game. game already exists"),
		},
		{
			description: "Start game without joining - fail",
			action: model.Action{
				Type:         model.ActionStart,
				GameID:       "ticTacToe",
				ActivePlayer: "Top",
			},
			conn: &MockConn{BytesWritten: make(chan []byte, 5)},
			games: map[model.GameID]*model.Game{"ticTacToe": {
				Id:          "ticTacToe",
				Connections: map[model.PlayerID]net.Conn{},
				Actions:     make(chan model.Action, 5),
			}},
			gameChan:    make(chan *model.Game, 2),
			expectedErr: fmt.Errorf("invalid action: start. not in a game"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			game, err := ConnectToGame(tc.action, tc.conn, tc.games, tc.gameChan)
			if err == nil {
				assert.Equal(t, tc.expectedGame.Id, game.Id)
				assert.Equal(t, len(tc.expectedGame.Connections), len(game.Connections))
				for player := range tc.expectedGame.Connections {
					_, ok := game.Connections[player]
					assert.True(t, ok, "Player '%s' expected in game.Connections", player)
				}
			}
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.conn.Closed, tc.expectedConnClosed)
		})
	}
}

func TestHandleNewGames(t *testing.T) {
	games := map[model.GameID]*model.Game{"go": {Id: "go"}}
	gameChan := make(chan *model.Game)
	go HandleNewGames(games, gameChan)

	gameChan <- &model.Game{Id: "chess"}
	goConn := &MockConn{}
	gameChan <- &model.Game{Id: "go", Connections: map[model.PlayerID]net.Conn{"top": goConn}}
	gameChan <- &model.Game{Id: "chess"} //flush

	assert.Equal(t, 2, len(games))
	assert.True(t, goConn.Closed)
}
