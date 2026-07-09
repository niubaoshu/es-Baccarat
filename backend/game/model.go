package game

import "github.com/niubaoshu/es-Baccarat/backend/card"

type Role string

const (
	Player       Role = "player"
	Banker       Role = "banker"
	Dealer       Role = "Dealer"
	TPPPPS       Role = "TPPPPS"
	PlayerDealer Role = "PlayerDealer"
)

type Person struct {
	name    string
	role    Role
	Balance int
}

type Table struct {
	Id         string
	Shoe       *card.Shoe
	Dealer     *Person
	TPPPPS     *Person
	Banker     *Person
	Players    []*Person
	minBet     int
	maxBet     int
	maxTie     int
	maxPanda8  int
	maxDragon7 int
	wagerRules []wagerRule
}

type wagerRule struct {
	min          int
	max          int
	Player       int
	PlayerDealer int
}

var wagerRules = map[int][]wagerRule{
	10: {
		{min: 10, max: 300, Player: 0, PlayerDealer: 2},
		{min: 300, max: 500, Player: 0, PlayerDealer: 4},
		{min: 500, max: 1500, Player: 0, PlayerDealer: 7},
		{min: 1500, max: 2500, Player: 0, PlayerDealer: 11},
		{min: 2500, max: 10000, Player: 0, PlayerDealer: 17},
	},
	50: {
		{min: 50, max: 300, Player: 0, PlayerDealer: 4},
		{min: 300, max: 1000, Player: 0, PlayerDealer: 12},
		{min: 1000, max: 2000, Player: 0, PlayerDealer: 17},
		{min: 2000, max: 5000, Player: 0, PlayerDealer: 22},
		{min: 5000, max: 50000, Player: 0, PlayerDealer: 52},
	},
}
