package rules

import "github.com/niubaoshu/es-Baccarat/backend/model"

// Outcome represents the final result of a Baccarat hand.
// Outcome 代表一局百家乐的最终结果。
type Outcome string

const (
	OutcomePlayer  Outcome = "Player"
	OutcomeBanker  Outcome = "Banker"
	OutcomeTie     Outcome = "Tie"
	OutcomeDragon7 Outcome = "Dragon 7"
	OutcomePanda8  Outcome = "Panda 8"
)

// DetermineOutcome compares final hands and returns the specific EZ Baccarat outcome.
// DetermineOutcome 比较最终手牌并返回具体的 EZ 百家乐结果。
func DetermineOutcome(playerHand *model.Hand, bankerHand *model.Hand) Outcome {
	pPts := playerHand.TotalPoints()
	bPts := bankerHand.TotalPoints()

	if pPts == bPts {
		return OutcomeTie
	}

	if pPts > bPts {
		// Player Wins. Check for Panda 8
		// 闲家赢。检查是否为熊猫8
		if pPts == 8 && len(playerHand.Cards) == 3 {
			return OutcomePanda8
		}
		return OutcomePlayer
	}

	// Banker Wins. Check for Dragon 7
	// 庄家赢。检查是否为龙7
	if bPts == 7 && len(bankerHand.Cards) == 3 {
		return OutcomeDragon7
	}
	return OutcomeBanker
}

// PayoutResult represents the result calculation for a single bet.
// PayoutResult 代表单笔投注的赔付计算结果。
type PayoutResult struct {
	WinAmount int // Net win amount (not including original bet if kept)
	// 净赢金额（不含已保留的原始投注）
	Returned int // Amount returned to player (e.g. original bet on Push or Win)
	// 退还给玩家的金额（例如平局或赢牌时退还的原始投注）
}

// NetChange returns the net change to the player's balance (WinAmount + Returned - OriginalBet)
// NetChange 返回玩家余额的净变动额（赢金额 + 退还金额 - 原始投注）
func (p PayoutResult) NetChange(originalBet int) int {
	return (p.WinAmount + p.Returned) - originalBet
}

// CalculatePayout takes an outcome, a bet type, and a bet amount,
// and returns the payout details (WinAmount and Returned amount).
// CalculatePayout 接收结果、下注类型和下注金额，
// 并返回赔付详情（赢金额和退还金额）。
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
			// 平局退注
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
			// 平局退注
			return PayoutResult{WinAmount: 0, Returned: betAmount}
		}
	}

	// Any other combo is a loss
	// 其他任何组合均视为输注
	return PayoutResult{WinAmount: 0, Returned: 0}
}
