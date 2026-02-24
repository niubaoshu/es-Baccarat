package rules

import (
	"testing"

	"github.com/niubaoshu/es-Baccarat/model"
)

func TestDetermineOutcome(t *testing.T) {
	tests := []struct {
		name       string
		playerHand *model.Hand
		bankerHand *model.Hand
		expected   Outcome
	}{
		{
			"Normal Player Win",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Nine}, {Suit: model.Spades, Rank: model.Jack}}},  // 9
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Eight}, {Suit: model.Spades, Rank: model.Jack}}}, // 8
			OutcomePlayer,
		},
		{
			"Normal Banker Win",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Spades, Rank: model.Three}}}, // 5
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Four}, {Suit: model.Spades, Rank: model.Two}}},  // 6
			OutcomeBanker,
		},
		{
			"Tie",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Seven}}},
			&model.Hand{Cards: []model.Card{{Suit: model.Hearts, Rank: model.Seven}}},
			OutcomeTie,
		},
		{
			"Dragon 7 (3 cards)",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Six}}},                                                                                 // 6
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Spades, Rank: model.Two}, {Suit: model.Spades, Rank: model.Three}}}, // 7 with 3 cards
			OutcomeDragon7,
		},
		{
			"Banker 7 with 2 cards is NOT Dragon",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Six}}},                                           // 6
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Three}, {Suit: model.Spades, Rank: model.Four}}}, // 7 with 2 cards
			OutcomeBanker,
		},
		{
			"Panda 8 (3 cards)",
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Three}, {Suit: model.Spades, Rank: model.Three}, {Suit: model.Spades, Rank: model.Two}}}, // 8 with 3 cards
			&model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Seven}}},                                                                                 // 7
			OutcomePanda8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineOutcome(tt.playerHand, tt.bankerHand)
			if got != tt.expected {
				t.Errorf("DetermineOutcome() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculatePayout(t *testing.T) {
	tests := []struct {
		name         string
		outcome      Outcome
		betType      BetType
		betAmount    int
		wantWin      int
		wantReturned int
		wantNet      int
	}{
		// Player Wins
		{"Player Bet on OutcomePlayer", OutcomePlayer, Player, 100, 100, 100, 100},
		{"Banker Bet on OutcomePlayer", OutcomePlayer, Banker, 100, 0, 0, -100},
		{"Panda Bet on OutcomePlayer", OutcomePlayer, Panda, 100, 0, 0, -100},

		// Banker Wins
		{"Banker Bet on OutcomeBanker", OutcomeBanker, Banker, 100, 100, 100, 100},
		{"Player Bet on OutcomeBanker", OutcomeBanker, Player, 100, 0, 0, -100},

		// Tie
		{"Tie Bet on OutcomeTie", OutcomeTie, Tie, 10, 80, 10, 80},
		{"Player Bet on OutcomeTie (Push)", OutcomeTie, Player, 100, 0, 100, 0},
		{"Banker Bet on OutcomeTie (Push)", OutcomeTie, Banker, 100, 0, 100, 0},
		{"Dragon Bet on OutcomeTie", OutcomeTie, Dragon, 100, 0, 0, -100},

		// Dragon 7
		{"Banker Bet on Dragon 7 (Push)", OutcomeDragon7, Banker, 100, 0, 100, 0},
		{"Dragon Bet on Dragon 7", OutcomeDragon7, Dragon, 10, 400, 10, 400},
		{"Player Bet on Dragon 7", OutcomeDragon7, Player, 100, 0, 0, -100},
		{"Tie Bet on Dragon 7", OutcomeDragon7, Tie, 100, 0, 0, -100},

		// Panda 8
		{"Player Bet on Panda 8", OutcomePanda8, Player, 100, 100, 100, 100},
		{"Panda Bet on Panda 8", OutcomePanda8, Panda, 10, 250, 10, 250},
		{"Banker Bet on Panda 8", OutcomePanda8, Banker, 100, 0, 0, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePayout(tt.outcome, tt.betType, tt.betAmount)
			if result.WinAmount != tt.wantWin {
				t.Errorf("WinAmount got %d, want %d", result.WinAmount, tt.wantWin)
			}
			if result.Returned != tt.wantReturned {
				t.Errorf("Returned got %d, want %d", result.Returned, tt.wantReturned)
			}
			if result.NetChange(tt.betAmount) != tt.wantNet {
				t.Errorf("NetChange got %d, want %d", result.NetChange(tt.betAmount), tt.wantNet)
			}
		})
	}
}
