package model

// Hand represents a player's or banker's hand of cards in Baccarat.
type Hand struct {
	Cards []Card
}

// AddCard adds a card to the hand.
func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

// TotalPoints returns the Baccarat point value of the hand (0-9).
// It sums the point values of all cards and drops the tens digit.
func (h *Hand) TotalPoints() int {
	total := 0
	for _, c := range h.Cards {
		total += c.PointValue()
	}
	return total % 10
}

// IsNatural checks if the hand is a "Natural" win (total of 8 or 9 in the first two cards).
func (h *Hand) IsNatural() bool {
	if len(h.Cards) == 2 {
		points := h.TotalPoints()
		if points == 8 || points == 9 {
			return true
		}
	}
	return false
}

// String provides a simple view of the cards in the hand and its total points.
func (h *Hand) String() string {
	res := ""
	for i, c := range h.Cards {
		if i > 0 {
			res += ", "
		}
		res += c.String()
	}
	return res
}
