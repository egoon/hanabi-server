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
