package model

import (
	"math/rand"
	"time"
)

type GameID string
type PlayerID string

type GameState struct {
	Id           GameID   `json:"id,omitempty"`
	Players      []Player `json:"players"`
	Clues        int      `json:"clues"`
	Lives        int      `json:"lives"`
	Discards     []Card   `json:"discards"`
	Table        []Card   `json:"table"`
	Deck         int      `json:"deck"`
	PlayedAction Action   `json:"playedAction"`
	Started      bool     `json:"started"`
	Ended        bool     `json:"ended"`
}

type Player struct {
	Id    PlayerID `json:"id"`
	Cards []Card   `json:"cards,omitempty"`
}

type Card struct {
	Color string `json:"color"`
	Value string `json:"value"`
}

func (g *GameState) ForPlayer(playerID PlayerID) (GameState, bool) {
	filtered := make([]Player, len(g.Players))
	ok := false
	for i, player := range g.Players {
		if player.Id == playerID {
			filtered[i] = Player{Id: playerID}
			ok = true
		} else {
			filtered[i] = player
		}
	}
	return GameState{
		Id:           g.Id,
		Players:      filtered,
		Clues:        g.Clues,
		Lives:        g.Lives,
		Discards:     g.Discards,
		Table:        g.Table,
		Deck:         g.Deck,
		PlayedAction: g.PlayedAction,
		Started:      g.Started,
		Ended:        g.Ended,
	}, ok
}

func (g *GameState) HasPlayer(player PlayerID) bool {
	for _, p := range g.Players {
		if p.Id == player {
			return true
		}
	}
	return false
}

func CreateDeck() []Card {
	colors := []string{"B", "G", "R", "W", "Y"}
	values := []string{"1", "1", "1", "2", "2", "3", "3", "4", "4", "5"}
	deck := make([]Card, len(colors)*len(values))
	i := 0
	for _, color := range colors {
		for _, value := range values {
			deck[i] = Card{
				Color: color,
				Value: value,
			}
		}
		i++
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}
