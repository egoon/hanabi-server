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
	Game         GameID   `json:"game,omitempty"`
	ActivePlayer PlayerID `json:"player,omitempty"`
	TargetPlayer PlayerID `json:"player,omitempty"`
	Card         []int    `json:"card,omitempty"`
	Clue         string   `json:"clue,omitempty"`
}
