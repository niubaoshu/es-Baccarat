package model

// Hand represents a player's or banker's hand of cards in Baccarat.
// Hand 代表百家乐中闲家或庄家的一手牌。
type Hand struct {
	Cards []Card
}

// AddCard adds a card to the hand.
// AddCard 向手牌中添加一张牌。
func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

// TotalPoints returns the Baccarat point value of the hand (0-9).
// It sums the point values of all cards and drops the tens digit.
// TotalPoints 返回该手牌的百家乐点数（0-9）。
// 它将所有牌的点数相加，并取个位数。
func (h *Hand) TotalPoints() int {
	total := 0
	for _, c := range h.Cards {
		total += c.PointValue()
	}
	return total % 10
}

// IsNatural checks if the hand is a "Natural" win (total of 8 or 9 in the first two cards).
// IsNatural 检查该手牌是否为"天牌"（前两张牌合计为 8 或 9）。
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
// String 提供手牌中所有牌及其总点数的简要视图。
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
