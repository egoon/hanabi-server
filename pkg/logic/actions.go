package logic

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/egoon/hanabi-server/pkg/model"
)

const (
	maxLives = 3
	maxClues = 8
)

func HandleGameActions(game *model.Game, deck []model.Card) {
	state := model.GameState{
		Id:           game.Id,
		Players:      make([]model.Player, 0, 5),
		Clues:        0,
		Lives:        0,
		Discards:     make([]model.Card, 0, len(deck)),
		Table:        make([]model.Card, 0, len(deck)/2),
		Deck:         0,
		PlayedAction: model.Action{},
		Started:      false,
		Ended:        false,
		Colors:       len(deck) / 10,
	}
	game.State = &state
	for {
		action := <-game.Actions
		state.PlayedAction = action
		deck = handleAction(action, &state, deck)
		sendStateToPlayers(&state, game.Connections)
		if state.Ended {
			for _, c := range game.Connections {
				c.Close()
			}
			log.Info("Game ", game.Id, " score: ", len(game.State.Table))
			break
		}
	}
}

func handleAction(action model.Action, state *model.GameState, deck []model.Card) []model.Card {
	switch action.Type {
	case model.ActionPing:
		// do nothing
	case model.ActionJoin:
		if len(state.Players) < 5 && !state.Started {
			state.Players = append(state.Players, model.Player{Id: action.ActivePlayer})
		}
	case model.ActionStart:
		cardsPerPlayer := 5
		if len(state.Players) > 3 {
			cardsPerPlayer = 4
		}
		for i := range state.Players {
			state.Players[i].Cards = deck[:cardsPerPlayer]
			deck = deck[cardsPerPlayer:]
		}
		state.Clues = maxClues
		state.Lives = maxLives
		state.Started = true
		state.Deck = len(deck)
	case model.ActionClue:
		state.Clues--
		for _, player := range state.Players {
			if player.Id == action.TargetPlayer {
				for i, card := range player.Cards {
					if card.Color == action.Clue || card.Value == action.Clue {
						action.Card = append(action.Card, i)
					}
				}
				break
			}
		}
	case model.ActionPlay:
		hand := state.Players[0].Cards
		card := hand[action.Card[0]]
		if isCardPlayable(card, state.Table) {
			state.Table = append(state.Table, card)
			if card.Value == "5" && state.Clues < maxClues {
				state.Clues++
			}
		} else {
			state.Discards = append(state.Discards, card)
			state.Lives--
		}
		hand[action.Card[0]], deck = drawCard(deck)
		if len(state.Table) == state.Colors*5 || state.Lives == 0 {
			state.Ended = true
		}
		state.Players[0].Cards = hand
	case model.ActionDiscard:
		hand := state.Players[0].Cards
		card := hand[action.Card[0]]
		state.Discards = append(state.Discards, card)
		hand[action.Card[0]], deck = drawCard(deck)
		state.Players[0].Cards = hand
		if state.Clues < maxClues {
			state.Clues++
		}
	}
	if len(deck) > 0 {
		state.Deck = len(deck)
	}
	if action.Type == model.ActionPlay || action.Type == model.ActionClue || action.Type == model.ActionDiscard {
		if len(deck) == 0 {
			state.Deck--
		}
		if -state.Deck == len(state.Players) {
			state.Ended = true
		}
		state.Players = append(state.Players[1:], state.Players[0])
	}
	state.PlayedAction = action
	return deck
}

func isCardPlayable(card model.Card, table []model.Card) bool {
	requiredCardPlayed := card.Value == "1"
	for _, c := range table {
		if c == card {
			return false
		}
		if c.Color == card.Color &&
			(c.Value == "1" && card.Value == "2" ||
				c.Value == "2" && card.Value == "3" ||
				c.Value == "3" && card.Value == "4" ||
				c.Value == "4" && card.Value == "5") {
			requiredCardPlayed = true
		}
	}
	return requiredCardPlayed
}

func ValidateAndCleanAction(action *model.Action, state *model.GameState) error {
	switch action.Type {
	case model.ActionPing:
		action.Card = nil
		action.Clue = ""
		action.GameID = ""
		action.TargetPlayer = ""
	case model.ActionCreate:
		if state != nil {
			return fmt.Errorf("already connected to a game")
		}
		action.Card = nil
		action.Clue = ""
		action.TargetPlayer = ""
	case model.ActionJoin:
		if state != nil {
			return fmt.Errorf("already connected to a game")
		}
		if action.GameID == "" {
			return fmt.Errorf("join action must have game id")
		}
		action.Card = nil
		action.Clue = ""
		action.TargetPlayer = ""
	case model.ActionStart:
		if state == nil {
			return fmt.Errorf("not connected to a game")
		}
		if state.Started {
			return fmt.Errorf("game already started")
		}
		if state.Players[0].Id != action.ActivePlayer {
			return fmt.Errorf("only creator may start game")
		}
		if len(state.Players) < 2 {
			return fmt.Errorf("too few players")
		}
		action.Card = nil
		action.Clue = ""
		action.GameID = ""
		action.TargetPlayer = ""
	case model.ActionClue:
		if state == nil {
			return fmt.Errorf("not connected to a game")
		}
		if !state.Started {
			return fmt.Errorf("game is not started")
		}
		if state.Players[0].Id != action.ActivePlayer {
			return fmt.Errorf("not your turn")
		}
		if state.Clues < 1 {
			return fmt.Errorf("there are no clues available to give")
		}
		pattern := `^[12345BGRWY]$`
		match, _ := regexp.MatchString(pattern, action.Clue)
		if !match {
			return fmt.Errorf("clue action must have clue field that matches %s", pattern)
		}
		if !state.HasPlayer(action.TargetPlayer) {
			return fmt.Errorf("player %s is not in this game", action.TargetPlayer)
		}
		if action.TargetPlayer == action.ActivePlayer {
			return fmt.Errorf("you may not target yourself")
		}
		action.GameID = ""
		action.Card = make([]int, 5)[:0]
	case model.ActionPlay:
		if state == nil {
			return fmt.Errorf("not connected to a game")
		}
		if !state.Started {
			return fmt.Errorf("game is not started")
		}
		if state.Players[0].Id != action.ActivePlayer {
			return fmt.Errorf("not your turn")
		}
		if len(action.Card) != 1 {
			return fmt.Errorf("exactly 1 card must be played. Not %d", len(action.Card))
		}
		if action.Card[0] < 0 || action.Card[0] >= len(state.Players[0].Cards) {
			return fmt.Errorf("no card on index %d", action.Card[0])
		}
		action.GameID = ""
		action.Clue = ""
		action.TargetPlayer = ""
	case model.ActionDiscard:
		if state == nil {
			return fmt.Errorf("not connected to a game")
		}
		if !state.Started {
			return fmt.Errorf("game is not started")
		}
		if state.Players[0].Id != action.ActivePlayer {
			return fmt.Errorf("not your turn")
		}
		if len(action.Card) != 1 {
			return fmt.Errorf("exactly 1 card must be discarded. Not %d", len(action.Card))
		}
		if action.Card[0] < 0 || action.Card[0] >= len(state.Players[0].Cards) {
			return fmt.Errorf("no card on index %d", action.Card[0])
		}
		action.GameID = ""
		action.Clue = ""
		action.TargetPlayer = ""
	default:
		return fmt.Errorf("unknown action: %s", action.Type)
	}
	return nil
}

func sendStateToPlayers(state *model.GameState, connections map[model.PlayerID]net.Conn) {
	for playerId, conn := range connections {
		if state.PlayedAction.Type == "ping" && state.PlayedAction.ActivePlayer != playerId {
			// only respond to player who pinged
			continue
		}
		playerState, _ := state.ForPlayer(playerId)
		msg, _ := json.Marshal(playerState)
		conn.Write(msg)
	}
}

func drawCard(cards []model.Card) (model.Card, []model.Card) {
	if len(cards) > 0 {
		card := cards[0]
		cards := cards[1:]
		return card, cards
	}
	return model.Card{Color: "-", Value: "-"}, cards
}
