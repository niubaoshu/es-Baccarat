package model

import "testing"

func TestCardPointValue(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		expected int
	}{
		{"Ace is 1", Card{Spades, Ace}, 1},
		{"Two is 2", Card{Hearts, Two}, 2},
		{"Nine is 9", Card{Diamonds, Nine}, 9},
		{"Ten is 0", Card{Clubs, Ten}, 0},
		{"Jack is 0", Card{Spades, Jack}, 0},
		{"Queen is 0", Card{Hearts, Queen}, 0},
		{"King is 0", Card{Diamonds, King}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.PointValue(); got != tt.expected {
				t.Errorf("Card.PointValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCardString(t *testing.T) {
	c := Card{Spades, Ace}
	if c.String() != "A♠" {
		t.Errorf("Expected A♠, got %s", c.String())
	}
	c2 := Card{Hearts, Ten}
	if c2.String() != "10♥" {
		t.Errorf("Expected 10♥, got %s", c2.String())
	}
}
