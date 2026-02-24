package model

import "testing"

func TestHandTotalPoints(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected int
	}{
		{"Single Ace", []Card{{Spades, Ace}}, 1},
		{"Two Nines", []Card{{Spades, Nine}, {Hearts, Nine}}, 8},                             // 18 -> 8
		{"Ten and Five", []Card{{Spades, Ten}, {Hearts, Five}}, 5},                           // 15 -> 5
		{"Five and King", []Card{{Spades, Five}, {Hearts, King}}, 5},                         // 5 + 0 -> 5
		{"Three Cards", []Card{{Spades, Two}, {Hearts, Three}, {Diamonds, Four}}, 9},         // 2+3+4 = 9
		{"Three Cards over 10", []Card{{Spades, Five}, {Hearts, Six}, {Diamonds, Seven}}, 8}, // 5+6+7 = 18 -> 8
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hand{Cards: tt.cards}
			if got := h.TotalPoints(); got != tt.expected {
				t.Errorf("Hand.TotalPoints() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHandIsNatural(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected bool
	}{
		{"Natural 8", []Card{{Spades, Four}, {Hearts, Four}}, true},
		{"Natural 9", []Card{{Spades, Four}, {Hearts, Five}}, true},
		{"Not Natural 7", []Card{{Spades, Three}, {Hearts, Four}}, false},
		{"Not Natural 10", []Card{{Spades, Five}, {Hearts, Five}}, false},
		{"3 Cards 8 is Not Natural", []Card{{Spades, Two}, {Hearts, Two}, {Diamonds, Four}}, false},
		{"1 Card 8 is Not Natural", []Card{{Spades, Eight}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hand{Cards: tt.cards}
			if got := h.IsNatural(); got != tt.expected {
				t.Errorf("Hand.IsNatural() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHandAddCard(t *testing.T) {
	h := Hand{}
	h.AddCard(Card{Spades, Ace})
	if len(h.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(h.Cards))
	}
	h.AddCard(Card{Hearts, Two})
	if len(h.Cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(h.Cards))
	}
}
