package rules

import (
	"github.com/niubaoshu/es-Baccarat/backend/model"
)

// DeterminePlayerHit decides if the Player should draw a third card.
// According to Baccarat rules, if either has a Natural (8 or 9), no one hits.
// If not natural, Player hits on 0-5, stands on 6-9.
// DeterminePlayerHit 决定闲家是否应该补第三张牌。
// 根据百家乐规则，若任意一方为天牌（8 或 9），则双方均不补牌。
// 若非天牌，闲家点数 0-5 时补牌，6-9 时不补。
func DeterminePlayerHit(playerHand *model.Hand, bankerHand *model.Hand) bool {
	if playerHand.IsNatural() || bankerHand.IsNatural() {
		return false
	}

	if playerHand.TotalPoints() <= 5 {
		return true
	}
	return false
}

// DetermineBankerHit decides if the Banker should draw a third card.
// It requires knowing whether the Player has already hit, and what specific card they drew.
// If either has a Natural (8 or 9), no one hits.
// DetermineBankerHit 决定庄家是否应该补第三张牌。
// 需要知道闲家是否已补牌，以及闲家所补的具体牌面。
// 若任意一方为天牌（8 或 9），则双方均不补牌。
func DetermineBankerHit(bankerHand *model.Hand, playerHand *model.Hand, playerHit bool, playerThirdCard *model.Card) bool {
	if playerHand.IsNatural() || bankerHand.IsNatural() {
		return false
	}

	bankerPts := bankerHand.TotalPoints()

	// If player did NOT hit (stood on 6 or 7)
	// 若闲家未补牌（停在 6 或 7）
	if !playerHit {
		// Banker hits on 0-5, stands on 6-7
		// 庄家点数 0-5 时补牌，6-7 时不补
		return bankerPts <= 5
	}

	// If Player DID hit, the Banker drawing depends on the third card drawn by the player
	// 若闲家已补牌，庄家是否补牌取决于闲家所补的第三张牌
	if playerThirdCard == nil { // Should not happen in proper flow if playerHit is true
		// 若 playerHit 为 true，正常流程中该情况不应出现
		return false
	}

	p3Pts := playerThirdCard.PointValue()

	switch bankerPts {
	case 0, 1, 2:
		return true // Always hit
		// 始终补牌
	case 3:
		// Banker hits on 3, unless player's third card was an 8
		// 庄家点数为 3 时补牌，除非闲家第三张牌为 8
		if p3Pts != 8 {
			return true
		}
		return false
	case 4:
		// Banker hits on 4 if player's third card is 2-7
		// 庄家点数为 4 且闲家第三张牌为 2-7 时补牌
		if p3Pts >= 2 && p3Pts <= 7 {
			return true
		}
		return false
	case 5:
		// Banker hits on 5 if player's third card is 4-7
		// 庄家点数为 5 且闲家第三张牌为 4-7 时补牌
		if p3Pts >= 4 && p3Pts <= 7 {
			return true
		}
		return false
	case 6:
		// Banker hits on 6 if player's third card is 6 or 7
		// 庄家点数为 6 且闲家第三张牌为 6 或 7 时补牌
		if p3Pts == 6 || p3Pts == 7 {
			return true
		}
		return false
	case 7:
		// Banker always stands on 7
		// 庄家点数为 7 时始终不补牌
		return false
	}

	return false
}
