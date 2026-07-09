package rules

// BetType represents the different betting options in EZ Baccarat Panda 8.
// BetType 代表 EZ 百家乐熊猫8中的各种下注类型。
type BetType string

const (
	Player BetType = "Player"
	Banker BetType = "Banker"
	Tie    BetType = "Tie"
	Dragon BetType = "Dragon 7"
	Panda  BetType = "Panda 8"
)
