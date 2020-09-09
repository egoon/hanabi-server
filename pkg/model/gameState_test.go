package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameState_ForPlayer(t *testing.T) {
	state := GameState{
		Id: "id",
		Players: []Player{
			{
				Id: "p1",
				Cards: []Card{
					{
						Color: "W",
						Value: "1",
					},
					{
						Color: "W",
						Value: "2",
					},
					{
						Color: "W",
						Value: "3",
					},
					{
						Color: "W",
						Value: "4",
					},
				},
			},
			{
				Id: "p2",
				Cards: []Card{
					{
						Color: "B",
						Value: "1",
					},
					{
						Color: "B",
						Value: "2",
					},
					{
						Color: "B",
						Value: "3",
					},
					{
						Color: "B",
						Value: "4",
					},
				},
			},
		},
		Clues: 3,
		Lives: 2,
		Discards: []Card{
			{
				Color: "G",
				Value: "1",
			},
		},
		Table: []Card{
			{
				Color: "Y",
				Value: "1",
			},
			{
				Color: "Y",
				Value: "2",
			},
		},
		Deck:         30,
		PlayedAction: Action{},
	}

	testCases := []struct {
		description string
		GameState
		playerID   PlayerID
		expectedOk bool
	}{
		{
			description: "first player",
			GameState:   state,
			playerID:    "p1",
			expectedOk:  true,
		},
		{
			description: "second/last player",
			GameState:   state,
			playerID:    "p2",
			expectedOk:  true,
		},
		{
			description: "missing player",
			GameState:   state,
			playerID:    "p3",
			expectedOk:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			newState, ok := tc.GameState.ForPlayer(tc.playerID)
			assert.Equal(t, tc.expectedOk, ok)
			if ok {
				assert.NotEqual(t, tc.GameState, newState)
				assert.Equal(t, tc.GameState.PlayedAction, newState.PlayedAction)
				assert.Equal(t, tc.GameState.Deck, newState.Deck)
				assert.Equal(t, tc.GameState.Table, newState.Table)
				assert.Equal(t, tc.GameState.Discards, newState.Discards)
				assert.Equal(t, tc.GameState.Lives, newState.Lives)
				assert.Equal(t, tc.GameState.Clues, newState.Clues)
				assert.Equal(t, tc.GameState.Id, newState.Id)
				for _, player := range newState.Players {
					if player.Id == tc.playerID {
						assert.Nil(t, player.Cards)
					} else {
						assert.Less(t, 0, len(player.Cards))
					}
				}
			} else {
				assert.Equal(t, tc.GameState, newState)
			}
		})
	}
}

func TestGameState_CreateDeck(t *testing.T) {
	deck := CreateDeck()
	assert.Equal(t, 50, len(deck))
	cards := map[Card]int{}
	for _, card := range deck {
		if count, ok := cards[card]; ok {
			cards[card] = count + 1
		} else {
			cards[card] = 1
		}
	}
	assert.Equal(t, 25, len(cards))
	colors := map[string]int{}
	for card, cardCount := range cards {
		if count, ok := colors[card.Color]; ok {
			colors[card.Color] = count + cardCount
		} else {
			colors[card.Color] = cardCount
		}
		if card.Value == "1" {
			assert.Equal(t, 3, cardCount)
		} else if card.Value == "5" {
			assert.Equal(t, 1, cardCount)
		} else if card.Value == "4" || card.Value == "3" || card.Value == "2" {
			assert.Equal(t, 2, cardCount)
		} else {
			assert.Fail(t, card.Value, "is not a valid card value")
		}
	}
	assert.Equal(t, 5, len(colors))
	for color, count := range colors {
		assert.Equal(t, 10, count, color, "should have 10 cards")
	}
}
