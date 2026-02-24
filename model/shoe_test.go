package model

import (
	"testing"
)

func TestNewShoe(t *testing.T) {
	shoe := NewShoe(8, 14)
	if len(shoe.Cards) != 8*52 {
		t.Errorf("Expected 416 cards for 8 decks, got %d", len(shoe.Cards))
	}
	if shoe.CardsLeft() != 416 {
		t.Errorf("Expected 416 cards left, got %d", shoe.CardsLeft())
	}
	if shoe.IsPastCutCard() {
		t.Errorf("Newly created shoe shouldn't be past cut card")
	}
}

func TestShoeDraw(t *testing.T) {
	shoe := NewShoe(1, 0)
	card, err := shoe.Draw()
	if err != nil {
		t.Fatalf("Unexpected error drawing card: %v", err)
	}
	// First card populated is typically Ace of Spades given our loop
	if card.Suit != Spades || card.Rank != Ace {
		t.Errorf("Expected first draw to be Ace of Spades, got %v", card)
	}
	if shoe.CardsLeft() != 51 {
		t.Errorf("Expected 51 cards left, got %d", shoe.CardsLeft())
	}
}

func TestShoeDrawEmpty(t *testing.T) {
	shoe := NewShoe(1, 0)
	for i := 0; i < 52; i++ {
		_, _ = shoe.Draw()
	}
	_, err := shoe.Draw()
	if err != ErrShoeEmpty {
		t.Errorf("Expected ErrShoeEmpty, got %v", err)
	}
}

func TestShoeIsPastCutCard(t *testing.T) {
	shoe := NewShoe(1, 14) // 52 cards, cut card at 14

	// Draw 38 cards to exactly hit 14 remaining
	for i := 0; i < 38; i++ {
		_, _ = shoe.Draw()
	}
	if !shoe.IsPastCutCard() {
		t.Errorf("Expected shoe to be past cut card with 14 cards left")
	}

	// 15 remaining
	shoe2 := NewShoe(1, 14)
	for i := 0; i < 37; i++ {
		_, _ = shoe2.Draw()
	}
	if shoe2.IsPastCutCard() {
		t.Errorf("Expected shoe to NOT be past cut card with 15 cards left")
	}
}

func TestShoeBurn(t *testing.T) {
	shoe := NewShoe(1, 0) // Unshuffled
	// 1st card is Ace of Spades (Rank 1). It dictates 1 card burned.
	// Total cards consumed from shoe = 1 (face-up) + 1 (burned) = 2.
	err := shoe.Burn()
	if err != nil {
		t.Fatalf("Unexpected error during burn: %v", err)
	}
	if shoe.CardsLeft() != 50 {
		t.Errorf("Expected 50 cards left after burning for Ace, got %d", shoe.CardsLeft())
	}

	// Let's create a custom shoe where the first card is a Jack
	shoe2 := NewShoe(1, 0)
	shoe2.Cards[0] = Card{Suit: Spades, Rank: Jack}
	// Jack is Rank 11, should burn 10 cards.
	// Total consumed = 1 (Jack) + 10 = 11.
	err = shoe2.Burn()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if shoe2.CardsLeft() != 41 {
		t.Errorf("Expected 41 cards left after burning for Jack, got %d", shoe2.CardsLeft())
	}
}
