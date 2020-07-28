package logic

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
)

var (
	b1     = model.Card{Color: "B", Value: "1"}
	b2     = model.Card{Color: "B", Value: "2"}
	b3     = model.Card{Color: "B", Value: "3"}
	b4     = model.Card{Color: "B", Value: "4"}
	b5     = model.Card{Color: "B", Value: "5"}
	g1     = model.Card{Color: "G", Value: "1"}
	g2     = model.Card{Color: "G", Value: "2"}
	g3     = model.Card{Color: "G", Value: "3"}
	g4     = model.Card{Color: "G", Value: "4"}
	g5     = model.Card{Color: "G", Value: "5"}
	r1     = model.Card{Color: "R", Value: "1"}
	r2     = model.Card{Color: "R", Value: "2"}
	r3     = model.Card{Color: "R", Value: "3"}
	r4     = model.Card{Color: "R", Value: "4"}
	r5     = model.Card{Color: "R", Value: "5"}
	w1     = model.Card{Color: "W", Value: "1"}
	w2     = model.Card{Color: "W", Value: "2"}
	w3     = model.Card{Color: "W", Value: "3"}
	w4     = model.Card{Color: "W", Value: "4"}
	w5     = model.Card{Color: "W", Value: "5"}
	y1     = model.Card{Color: "Y", Value: "1"}
	y2     = model.Card{Color: "Y", Value: "2"}
	y3     = model.Card{Color: "Y", Value: "3"}
	y4     = model.Card{Color: "Y", Value: "4"}
	y5     = model.Card{Color: "Y", Value: "5"}
	noCard = model.Card{Color: "-", Value: "-"}
)

func TestValidateAndCleanAction(t *testing.T) {
	testCases := []struct {
		description    string
		action         model.Action
		state          *model.GameState
		expectedAction model.Action
		expectedError  error
	}{
		{
			description:    "Clean Ping - OK",
			action:         model.Action{Type: model.ActionPing},
			state:          &model.GameState{},
			expectedAction: model.Action{Type: model.ActionPing},
			expectedError:  nil,
		},
		//PING
		{
			description: "Dirty Ping - OK",
			action: model.Action{
				Type:         model.ActionPing,
				Game:         "Dirty",
				ActivePlayer: "Active Player",
				TargetPlayer: "Dirty",
				Card:         []int{1, 2, 3},
				Clue:         "Dirty",
			},
			state: &model.GameState{},
			expectedAction: model.Action{
				Type:         model.ActionPing,
				ActivePlayer: "Active Player",
			},
			expectedError: nil,
		},
		{
			description: "Clean Create - OK",
			action: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: nil,
		},
		//CREATE
		{
			description: "Dirty Create - OK",
			action: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "Dirty",
				Card:         []int{1, 2, 3},
				Clue:         "Dirty",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: nil,
		},
		{
			description: "Clean Create - Fail: game already created",
			action: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			state: &model.GameState{Id: "Already Created"},
			expectedAction: model.Action{
				Type:         model.ActionCreate,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: fmt.Errorf("already connected to a game"),
		},
		//JOIN
		{
			description: "Clean Join - OK",
			action: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: nil,
		},
		{
			description: "Dirty Join - OK",
			action: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "Dirty",
				Card:         []int{0, 1, 2},
				Clue:         "Dirty",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: nil,
		},
		{
			description: "Clean Join - Fail: already connected",
			action: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			state: &model.GameState{},
			expectedAction: model.Action{
				Type:         model.ActionJoin,
				Game:         "My Game",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: fmt.Errorf("already connected to a game"),
		},
		{
			description: "Clean Join - Fail: missing game id",
			action: model.Action{
				Type:         model.ActionJoin,
				Game:         "",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionJoin,
				Game:         "",
				ActivePlayer: "Active Player",
				TargetPlayer: "",
				Card:         nil,
				Clue:         "",
			},
			expectedError: fmt.Errorf("join action must have game id"),
		},
		//START
		{
			description: "Clean Start - OK",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
			},
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: nil,
		},
		{
			description: "Dirty Start - OK",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
				Game:         "Dirty",
				Card:         nil,
				Clue:         "Dirty",
				TargetPlayer: "Dirty",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
			},
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: nil,
		},
		{
			description: "Clean Start - Fail: no game",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: fmt.Errorf("not connected to a game"),
		},
		{
			description: "Clean Start - Fail: game already started",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: fmt.Errorf("game already started"),
		},
		{
			description: "Clean Start - Fail: not first player",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "You"}, {Id: "Me"}},
			},
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: fmt.Errorf("only creator may start game"),
		},
		{
			description: "Clean Start - Fail: too few players",
			action: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}},
			},
			expectedAction: model.Action{
				Type:         model.ActionStart,
				ActivePlayer: "Me",
			},
			expectedError: fmt.Errorf("too few players"),
		},
		//CLUE
		{
			description: "Clean Clue Blue - OK",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
				Card:         []int{},
			},
			expectedError: nil,
		},
		{
			description: "Clean Clue 3 - OK",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "3",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "3",
				Card:         []int{},
			},
			expectedError: nil,
		},
		{
			description: "Dirty Clue White - OK",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "W",
				Card:         []int{1, 2, 3, 4, 5},
				Game:         "Play first card!!",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "W",
				Card:         []int{},
			},
			expectedError: nil,
		}, {
			description: "Clean Clue Blue - Fail: no game",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("not connected to a game"),
		},
		{
			description: "Clean Clue Blue - Fail: game not started",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: false,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("game is not started"),
		},
		{
			description: "Clean Clue Blue - Fail: not your turn",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "You"}, {Id: "Me"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("not your turn"),
		},
		{
			description: "Clean Clue Purple - Fail: purple is not a valid color",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "P",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "P",
			},
			expectedError: fmt.Errorf("clue action must have clue field that matches ^[12345BGRWY]$"),
		},
		{
			description: "Clean Clue 6 - Fail: 6 is not a valid value",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "6",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "6",
			},
			expectedError: fmt.Errorf("clue action must have clue field that matches ^[12345BGRWY]$"),
		},
		{
			description: "Clean Clue Blue - Fail: target player not in game",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "Them",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "Them",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("player Them is not in this game"),
		},
		{
			description: "Clean Clue Blue - Fail: can't give clues to self",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "Me",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   8,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "Me",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("you may not target yourself"),
		},
		{
			description: "Clean Clue Blue - fail: no clues left",
			action: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me"}, {Id: "You"}},
				Started: true,
				Clues:   0,
			},
			expectedAction: model.Action{
				Type:         model.ActionClue,
				ActivePlayer: "Me",
				TargetPlayer: "You",
				Clue:         "B",
			},
			expectedError: fmt.Errorf("there are no clues available to give"),
		},
		//PLAY
		{
			description: "Clean Play first card - OK",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: nil,
		},
		{
			description: "Clean Play last card - OK",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{4},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{4},
			},
			expectedError: nil,
		},
		{
			description: "Dirty Play middle card - OK",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{2},
				TargetPlayer: "Dirty",
				Game:         "Dirty",
				Clue:         "Dirty",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{2},
			},
			expectedError: nil,
		},
		{
			description: "Clean Play first card - Fail: no game",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("not connected to a game"),
		},
		{
			description: "Clean Play first card - Fail: not your turn",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "You"}, {Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("not your turn"),
		},
		{
			description: "Clean Play first card - Fail: not started",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: false,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("game is not started"),
		},
		{
			description: "Clean Play first card - Fail: play 0 cards",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{},
			},
			expectedError: fmt.Errorf("exactly 1 card must be played. Not 0"),
		},
		{
			description: "Clean Play first card - Fail: play 2 cards",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{2, 3},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{2, 3},
			},
			expectedError: fmt.Errorf("exactly 1 card must be played. Not 2"),
		},
		{
			description: "Clean Play first card - Fail: incorrect index",
			action: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{5},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionPlay,
				ActivePlayer: "Me",
				Card:         []int{5},
			},
			expectedError: fmt.Errorf("no card on index 5"),
		},
		//DISCARD
		{
			description: "Clean Discard first card - OK",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: nil,
		},
		{
			description: "Clean Discard last card - OK",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{4},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{4},
			},
			expectedError: nil,
		},
		{
			description: "Dirty Discard middle card - OK",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{2},
				TargetPlayer: "Dirty",
				Game:         "Dirty",
				Clue:         "Dirty",
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{2},
			},
			expectedError: nil,
		},
		{
			description: "Clean Discard first card - Fail: no game",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: nil,
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("not connected to a game"),
		},
		{
			description: "Clean Discard first card - Fail: not your turn",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "You"}, {Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("not your turn"),
		},
		{
			description: "Clean Discard first card - Fail: not started",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: false,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{0},
			},
			expectedError: fmt.Errorf("game is not started"),
		},
		{
			description: "Clean Discard first card - Fail: Discard 0 cards",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{},
			},
			expectedError: fmt.Errorf("exactly 1 card must be discarded. Not 0"),
		},
		{
			description: "Clean Discard first card - Fail: Discard 2 cards",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{2, 3},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{2, 3},
			},
			expectedError: fmt.Errorf("exactly 1 card must be discarded. Not 2"),
		},
		{
			description: "Clean Discard first card - Fail: incorrect index",
			action: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{5},
			},
			state: &model.GameState{
				Players: []model.Player{{Id: "Me", Cards: []model.Card{w1, w1, w2, w2, w3}}, {Id: "You"}},
				Started: true,
			},
			expectedAction: model.Action{
				Type:         model.ActionDiscard,
				ActivePlayer: "Me",
				Card:         []int{5},
			},
			expectedError: fmt.Errorf("no card on index 5"),
		},
		{
			description: "Unknown action - Fail",
			action: model.Action{
				Type: "unknown action",
			},
			expectedAction: model.Action{
				Type: "unknown action",
			},
			expectedError: fmt.Errorf("unknown action: unknown action"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := ValidateAndCleanAction(&tc.action, tc.state)
			assert.Equal(t, tc.expectedAction, tc.action)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestHandleAction(t *testing.T) {
	testCases := []struct {
		description   string
		action        model.Action
		state         model.GameState
		deck          []model.Card
		expectedState model.GameState
		expectedDeck  []model.Card
	}{
		//PING
		{
			description:   "Ping - OK",
			action:        model.Action{Type: model.ActionPing},
			state:         model.GameState{},
			expectedState: model.GameState{PlayedAction: model.Action{Type: model.ActionPing}},
		},
		//JOIN
		{
			description: "Join empty game",
			action:      model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			state:       model.GameState{},
			expectedState: model.GameState{
				Players:      []model.Player{{Id: "Up"}},
				PlayedAction: model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			},
		},
		{
			description: "Join 2 player game",
			action:      model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			state:       model.GameState{Players: []model.Player{{Id: "Down"}, {Id: "Strange"}}},
			expectedState: model.GameState{
				Players:      []model.Player{{Id: "Down"}, {Id: "Strange"}, {Id: "Up"}},
				PlayedAction: model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			},
		},
		{
			description: "Join 5 player game - fail",
			action:      model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			state:       model.GameState{Players: []model.Player{{Id: "Down"}, {Id: "Strange"}, {Id: "Charm"}, {Id: "Top"}, {Id: "Bottom"}}},
			expectedState: model.GameState{Players: []model.Player{
				{Id: "Down"}, {Id: "Strange"}, {Id: "Charm"}, {Id: "Top"}, {Id: "Bottom"}},
				PlayedAction: model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			},
		},
		{
			description: "Join started game - fail",
			action:      model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			state:       model.GameState{Players: []model.Player{{Id: "Down"}, {Id: "Strange"}}, Started: true},
			expectedState: model.GameState{
				Players:      []model.Player{{Id: "Down"}, {Id: "Strange"}},
				Started:      true,
				PlayedAction: model.Action{Type: model.ActionJoin, ActivePlayer: "Up"},
			},
		},
		//START
		{
			description: "Start 2 player game",
			action:      model.Action{Type: model.ActionStart, ActivePlayer: "Up"},
			state:       model.GameState{Players: []model.Player{{Id: "Up"}, {Id: "Down"}}},
			deck:        []model.Card{w1, w2, w3, w4, w5, r1, r2, r3, r4, r5, b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}}, //dealt 5 cards
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Started:      true,
				Deck:         5,
				Clues:        8,
				Lives:        3,
				PlayedAction: model.Action{Type: model.ActionStart, ActivePlayer: "Up"},
			},
			expectedDeck: []model.Card{b1, b2, b3, b4, b5}, //five last cards of the deck
		},
		{
			description: "Start 4 player game",
			action:      model.Action{Type: model.ActionStart, ActivePlayer: "Up"},
			state:       model.GameState{Players: []model.Player{{Id: "Up"}, {Id: "Down"}, {Id: "Strange"}, {Id: "Charm"}}},
			deck:        []model.Card{w1, w2, w3, w4, r1, r2, r3, r4, b1, b2, b3, b4, g1, g2, g3, g4},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4}}, //dealt 4 cards
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4}},
					{Id: "Strange", Cards: []model.Card{b1, b2, b3, b4}},
					{Id: "Charm", Cards: []model.Card{g1, g2, g3, g4}},
				},
				Started:      true,
				Deck:         0,
				Clues:        8,
				Lives:        3,
				PlayedAction: model.Action{Type: model.ActionStart, ActivePlayer: "Up"},
			},
			expectedDeck: []model.Card{},
		},
		//CLUE
		{
			description: "Clue red match all",
			action:      model.Action{Type: model.ActionClue, ActivePlayer: "Up", TargetPlayer: "Down", Clue: "R"},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Clues: 8,
				Deck:  5,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
				},
				Clues: 7, //reduced by one
				PlayedAction: model.Action{
					Type:         model.ActionClue,
					ActivePlayer: "Up",
					TargetPlayer: "Down",
					Clue:         "R",
					Card:         []int{0, 1, 2, 3, 4}, //cards added to action
				},
				Deck: 5,
			},
			expectedDeck: []model.Card{b1, b2, b3, b4, b5},
		},
		{
			description: "Clue white match none",
			action:      model.Action{Type: model.ActionClue, ActivePlayer: "Up", TargetPlayer: "Down", Clue: "W"},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Clues: 8,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
				},
				Clues: 7, //reduced by one
				Deck:  5,
				PlayedAction: model.Action{
					Type:         model.ActionClue,
					ActivePlayer: "Up",
					TargetPlayer: "Down",
					Clue:         "W",
					Card:         nil, //no cards added to action
				},
			},
			expectedDeck: []model.Card{b1, b2, b3, b4, b5},
		},
		{
			description: "Clue 3 match index 2",
			action:      model.Action{Type: model.ActionClue, ActivePlayer: "Up", TargetPlayer: "Down", Clue: "3"},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Clues: 8,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
				},
				Clues: 7, //reduced by one
				Deck:  5,
				PlayedAction: model.Action{
					Type:         model.ActionClue,
					ActivePlayer: "Up",
					TargetPlayer: "Down",
					Clue:         "3",
					Card:         []int{2}, //middle card added to action
				},
			},
			expectedDeck: []model.Card{b1, b2, b3, b4, b5},
		},
		//PLAY
		{
			description: "Play W1 - ok",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{b1, w2, w3, w4, w5}}, // index 0: w1 -> b1
				},
				PlayedAction: model.Action{
					Type:         model.ActionPlay,
					ActivePlayer: "Up",
					Card:         []int{0}, //cards added to action
				},
				Deck:  4, //reduced by one
				Lives: 3,
				Table: []model.Card{w1},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W1 - already played",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Table: []model.Card{w1},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{b1, w2, w3, w4, w5}}, // index 0: w1 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{0}},
				Deck:         4, //reduced by one
				Lives:        2, //reduced by one
				Table:        []model.Card{w1},
				Discards:     []model.Card{w1}, // played card ends up in discard
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W2 - ok",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{1}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Table: []model.Card{w1},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, b1, w3, w4, w5}}, // index 1: w2 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{1}},
				Deck:         4, //reduced by one
				Lives:        3,
				Table:        []model.Card{w1, w2},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W2 - fail: W1 not played",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{1}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Table: []model.Card{b1},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, b1, w3, w4, w5}}, // index 1: w2 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{1}},
				Deck:         4, //reduced by one
				Lives:        2, //reduced by one
				Table:        []model.Card{b1},
				Discards:     []model.Card{w2},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W5 - ok: extra clue",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Table: []model.Card{w1, w2, w3, w4},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, b1}}, // index 4: w5 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
				Deck:         4, //reduced by one
				Lives:        3,
				Clues:        1, //increased by one
				Table:        []model.Card{w1, w2, w3, w4, w5},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W5 - ok: already max clues",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Clues: 8,
				Table: []model.Card{w1, w2, w3, w4},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, b1}}, // index 4: w5 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
				Deck:         4, //reduced by one
				Lives:        3,
				Clues:        8, //not changed
				Table:        []model.Card{w1, w2, w3, w4, w5},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W5 - ok: end game",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Table: []model.Card{
					b1, b2, b3, b4, b5,
					g1, g2, g3, g4, g5,
					r1, r2, r3, r4, r5,
					y1, y2, y3, y4, y5,
					w1, w2, w3, w4},
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, b1}}, // index 4: w5 -> b1
				},
				PlayedAction: model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{4}},
				Deck:         4, //reduced by one
				Lives:        3,
				Clues:        1, //increased by one
				Table: []model.Card{
					b1, b2, b3, b4, b5,
					g1, g2, g3, g4, g5,
					r1, r2, r3, r4, r5,
					y1, y2, y3, y4, y5,
					w1, w2, w3, w4, w5},
				Ended: true,
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Play W1 empty deck",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  0,
				Lives: 3,
			},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{noCard, w2, w3, w4, w5}}, // index 0: w1 -> nil
				},
				PlayedAction: model.Action{
					Type:         model.ActionPlay,
					ActivePlayer: "Up",
					Card:         []int{0}, //cards added to action
				},
				Deck:  -1, //reduced by one
				Lives: 3,
				Table: []model.Card{w1},
			},
		},
		{
			description: "Play last card",
			action:      model.Action{Type: model.ActionPlay, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  -1,
				Lives: 3,
			},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{noCard, w2, w3, w4, w5}}, // index 0: w1 -> nil
				},
				PlayedAction: model.Action{
					Type:         model.ActionPlay,
					ActivePlayer: "Up",
					Card:         []int{0}, //cards added to action
				},
				Deck:  -2, //reduced by one
				Lives: 3,
				Table: []model.Card{w1},
				Ended: true,
			},
		},
		//DISCARD
		{
			description: "Discard W1",
			action:      model.Action{Type: model.ActionDiscard, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{b1, w2, w3, w4, w5}}, // index 0: w1 -> b1
				},
				PlayedAction: model.Action{
					Type:         model.ActionDiscard,
					ActivePlayer: "Up",
					Card:         []int{0}, //cards added to action
				},
				Deck:     4, //reduced by one
				Lives:    3,
				Clues:    1, //increased by one
				Discards: []model.Card{w1},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
		{
			description: "Discard W1 - no clue",
			action:      model.Action{Type: model.ActionDiscard, ActivePlayer: "Up", Card: []int{0}},
			state: model.GameState{
				Players: []model.Player{
					{Id: "Up", Cards: []model.Card{w1, w2, w3, w4, w5}},
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
				},
				Deck:  5,
				Lives: 3,
				Clues: 8,
			},
			deck: []model.Card{b1, b2, b3, b4, b5},
			expectedState: model.GameState{
				Players: []model.Player{
					{Id: "Down", Cards: []model.Card{r1, r2, r3, r4, r5}},
					{Id: "Up", Cards: []model.Card{b1, w2, w3, w4, w5}}, // index 0: w1 -> b1
				},
				PlayedAction: model.Action{
					Type:         model.ActionDiscard,
					ActivePlayer: "Up",
					Card:         []int{0}, //cards added to action
				},
				Deck:     4, //reduced by one
				Lives:    3,
				Clues:    8, //did not increase
				Discards: []model.Card{w1},
			},
			expectedDeck: []model.Card{b2, b3, b4, b5}, // first card removed
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			state, deck := handleAction(tc.action, tc.state, tc.deck)
			assert.Equal(t, tc.expectedState, state)
			assert.Equal(t, tc.expectedDeck, deck)
		})
	}
}

func TestHandleGameActions(t *testing.T) {
	deck := []model.Card{
		w1, w2, w3, w4, w5,
		b1, b1, b2, b3, b4,
		w1, b1, w2, b2, b5,
		w1, w3, b3, w4, b4,
	}
	actions := make(chan model.Action, 5)
	strangeConn := &model.MockConn{BytesWritten: make(chan []byte)}
	charmConn := &model.MockConn{BytesWritten: make(chan []byte)}
	game := model.Game{
		Id: "game",
		Connections: map[model.PlayerID]net.Conn{
			"Strange": strangeConn,
			"Charm":   charmConn,
		},
		Actions: actions,
		State:   nil,
	}
	turns := []struct {
		action        model.Action
		expectedState model.GameState
	}{
		{
			action: model.Action{Type: model.ActionJoin, ActivePlayer: "Strange"},
			expectedState: model.GameState{
				Id:           "game",
				Players:      []model.Player{{Id: "Strange"}},
				Clues:        0,
				Lives:        0,
				Discards:     []model.Card{},
				Table:        []model.Card{},
				Deck:         20,
				PlayedAction: model.Action{Type: model.ActionJoin},
				Started:      false,
				Ended:        false,
			},
		},
		{
			action: model.Action{Type: model.ActionJoin, ActivePlayer: "Charm"},
			expectedState: model.GameState{
				Id:           "game",
				Players:      []model.Player{{Id: "Strange"}, {Id: "Charm"}},
				Clues:        0,
				Lives:        0,
				Discards:     []model.Card{},
				Table:        []model.Card{},
				Deck:         20,
				PlayedAction: model.Action{Type: model.ActionJoin},
				Started:      false,
				Ended:        false,
			},
		},
	}
	go HandleGameActions(&game, deck)
	go func() {
		<-charmConn.BytesWritten
		log.Info("charm got state")
	}()
	for _, turn := range turns {
		log.Info("play action")
		actions <- turn.action
		log.Info("read response")
		strangeBytes := <-strangeConn.BytesWritten
		state := model.GameState{}
		err := json.Unmarshal(strangeBytes, &state)
		log.Info(state)
		assert.Nil(t, err, "failed to unmarshal state: ", strangeBytes)
		assert.Equal(t, turn.expectedState, state)
	}
}
