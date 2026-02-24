package rules

import (
	"github.com/niubaoshu/es-Baccarat/backend/model"
)

// DeterminePlayerHit decides if the Player should draw a third card.
// According to Baccarat rules, if either has a Natural (8 or 9), no one hits.
// If not natural, Player hits on 0-5, stands on 6-9.
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
func DetermineBankerHit(bankerHand *model.Hand, playerHand *model.Hand, playerHit bool, playerThirdCard *model.Card) bool {
	if playerHand.IsNatural() || bankerHand.IsNatural() {
		return false
	}

	bankerPts := bankerHand.TotalPoints()

	// If player did NOT hit (stood on 6 or 7)
	if !playerHit {
		// Banker hits on 0-5, stands on 6-7
		return bankerPts <= 5
	}

	// If Player DID hit, the Banker drawing depends on the third card drawn by the player
	if playerThirdCard == nil { // Should not happen in proper flow if playerHit is true
		return false
	}

	p3Pts := playerThirdCard.PointValue()

	switch bankerPts {
	case 0, 1, 2:
		return true // Always hit
	case 3:
		// Banker hits on 3, unless player's third card was an 8
		if p3Pts != 8 {
			return true
		}
		return false
	case 4:
		// Banker hits on 4 if player's third card is 2-7
		if p3Pts >= 2 && p3Pts <= 7 {
			return true
		}
		return false
	case 5:
		// Banker hits on 5 if player's third card is 4-7
		if p3Pts >= 4 && p3Pts <= 7 {
			return true
		}
		return false
	case 6:
		// Banker hits on 6 if player's third card is 6 or 7
		if p3Pts == 6 || p3Pts == 7 {
			return true
		}
		return false
	case 7:
		// Banker always stands on 7
		return false
	}

	return false
}
