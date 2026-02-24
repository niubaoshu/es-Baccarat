package rules

import (
	"testing"

	"github.com/niubaoshu/es-Baccarat/backend/model"
)

func TestDeterminePlayerHit(t *testing.T) {
	// Natural test
	pHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Nine}, {Suit: model.Hearts, Rank: model.Ten}}}
	bHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Hearts, Rank: model.Two}}}
	if DeterminePlayerHit(pHand, bHand) {
		t.Errorf("Player should not hit on a natural 9")
	}

	// Hit on 0-5
	pHand = &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Hearts, Rank: model.Three}}} // 5
	if !DeterminePlayerHit(pHand, bHand) {
		t.Errorf("Player should hit on a total of 5")
	}

	// Stand on 6-7
	pHand = &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Hearts, Rank: model.Four}}} // 6
	if DeterminePlayerHit(pHand, bHand) {
		t.Errorf("Player should stand on a total of 6")
	}
}

func TestDetermineBankerHit_PlayerStood(t *testing.T) {
	// Test rules when player stood
	var hitTests = []struct {
		bankerPts int
		expected  bool
	}{
		{0, true}, {1, true}, {2, true}, {3, true}, {4, true}, {5, true},
		{6, false}, {7, false},
	}

	for _, tt := range hitTests {
		bHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Rank(tt.bankerPts)}}}
		// For 0 pts, use a 10 rank card
		if tt.bankerPts == 0 {
			bHand = &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Ten}}}
		}

		pHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Hearts, Rank: model.Four}}} // 6, stood

		got := DetermineBankerHit(bHand, pHand, false, nil)
		if got != tt.expected {
			t.Errorf("Banker points %d (player stood): got hit=%v, want %v", tt.bankerPts, got, tt.expected)
		}
	}
}

func TestDetermineBankerHit_PlayerHit(t *testing.T) {
	// Let's test specific sub-rules when player hit

	// Example: Banker has 4. Banker hits if player's 3rd card is 2-7.
	bHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Four}}}
	pHand := &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Two}, {Suit: model.Hearts, Rank: model.Two}}}

	p3 := &model.Card{Suit: model.Hearts, Rank: model.Five} // point value 5
	if !DetermineBankerHit(bHand, pHand, true, p3) {
		t.Errorf("Banker (4) should hit if player 3rd card is 5")
	}

	p3 = &model.Card{Suit: model.Hearts, Rank: model.Eight} // point value 8
	if DetermineBankerHit(bHand, pHand, true, p3) {
		t.Errorf("Banker (4) should stand if player 3rd card is 8")
	}

	// Example: Banker has 6. Banker hits if player's 3rd card is 6 or 7.
	bHand = &model.Hand{Cards: []model.Card{{Suit: model.Spades, Rank: model.Six}}}

	p3 = &model.Card{Suit: model.Hearts, Rank: model.Seven} // point value 7
	if !DetermineBankerHit(bHand, pHand, true, p3) {
		t.Errorf("Banker (6) should hit if player 3rd card is 7")
	}

	p3 = &model.Card{Suit: model.Hearts, Rank: model.Five} // point value 5
	if DetermineBankerHit(bHand, pHand, true, p3) {
		t.Errorf("Banker (6) should stand if player 3rd card is 5")
	}
}
