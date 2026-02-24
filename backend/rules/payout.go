package rules

import "github.com/niubaoshu/es-Baccarat/backend/model"

// Outcome represents the final result of a Baccarat hand.
type Outcome string

const (
	OutcomePlayer  Outcome = "Player"
	OutcomeBanker  Outcome = "Banker"
	OutcomeTie     Outcome = "Tie"
	OutcomeDragon7 Outcome = "Dragon 7"
	OutcomePanda8  Outcome = "Panda 8"
)

// DetermineOutcome compares final hands and returns the specific EZ Baccarat outcome.
func DetermineOutcome(playerHand *model.Hand, bankerHand *model.Hand) Outcome {
	pPts := playerHand.TotalPoints()
	bPts := bankerHand.TotalPoints()

	if pPts == bPts {
		return OutcomeTie
	}

	if pPts > bPts {
		// Player Wins. Check for Panda 8
		if pPts == 8 && len(playerHand.Cards) == 3 {
			return OutcomePanda8
		}
		return OutcomePlayer
	}

	// Banker Wins. Check for Dragon 7
	if bPts == 7 && len(bankerHand.Cards) == 3 {
		return OutcomeDragon7
	}
	return OutcomeBanker
}

// PayoutResult represents the result calculation for a single bet.
type PayoutResult struct {
	WinAmount int // Net win amount (not including original bet if kept)
	Returned  int // Amount returned to player (e.g. original bet on Push or Win)
}

// NetChange returns the net change to the player's balance (WinAmount + Returned - OriginalBet)
func (p PayoutResult) NetChange(originalBet int) int {
	return (p.WinAmount + p.Returned) - originalBet
}

// CalculatePayout takes an outcome, a bet type, and a bet amount,
// and returns the payout details (WinAmount and Returned amount).
func CalculatePayout(outcome Outcome, betType BetType, betAmount int) PayoutResult {
	switch outcome {
	case OutcomePlayer:
		switch betType {
		case Player:
			return PayoutResult{WinAmount: betAmount, Returned: betAmount}
		}

	case OutcomePanda8:
		switch betType {
		case Player:
			return PayoutResult{WinAmount: betAmount, Returned: betAmount}
		case Panda:
			return PayoutResult{WinAmount: betAmount * 25, Returned: betAmount}
		}

	case OutcomeBanker:
		switch betType {
		case Banker:
			return PayoutResult{WinAmount: betAmount, Returned: betAmount}
		}

	case OutcomeDragon7:
		switch betType {
		case Banker:
			// Push
			return PayoutResult{WinAmount: 0, Returned: betAmount}
		case Dragon:
			return PayoutResult{WinAmount: betAmount * 40, Returned: betAmount}
		}

	case OutcomeTie:
		switch betType {
		case Tie:
			return PayoutResult{WinAmount: betAmount * 8, Returned: betAmount}
		case Player, Banker:
			// Push
			return PayoutResult{WinAmount: 0, Returned: betAmount}
		}
	}

	// Any other combo is a loss
	return PayoutResult{WinAmount: 0, Returned: 0}
}
