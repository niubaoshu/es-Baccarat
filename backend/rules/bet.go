package rules

// BetType represents the different betting options in EZ Baccarat Panda 8.
type BetType string

const (
	Player BetType = "Player"
	Banker BetType = "Banker"
	Tie    BetType = "Tie"
	Dragon BetType = "Dragon 7"
	Panda  BetType = "Panda 8"
)
