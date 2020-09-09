package model

const (
	ActionPing    = "ping"
	ActionCreate  = "create"
	ActionJoin    = "join"
	ActionStart   = "start"
	ActionClue    = "clue"
	ActionPlay    = "play"
	ActionDiscard = "discard"
)

type Action struct {
	Type         string   `json:"type"`
	GameID       GameID   `json:"game,omitempty"`
	ActivePlayer PlayerID `json:"activePlayer,omitempty"`
	TargetPlayer PlayerID `json:"targetPlayer,omitempty"`
	Card         []int    `json:"card,omitempty"`
	Clue         string   `json:"clue,omitempty"`
}
