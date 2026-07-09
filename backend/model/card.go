package model

// Suit represents the suit of a playing card.
// Suit 代表一张扑克牌的花色。
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// Rank represents the rank (face value) of a playing card.
// Rank 代表一张扑克牌的点数（面值）。
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card represents a standard playing card without a Joker.
// Card 代表一张不含大小王的标准扑克牌。
type Card struct {
	Suit Suit
	Rank Rank
}

// PointValue calculates the Baccarat point value of the card.
// 2-9 retain their value. 10, J, Q, K are worth 0. Ace is worth 1.
// PointValue 计算该牌在百家乐中的点数值。
// 2-9 保留其原始点数，10、J、Q、K 点数为 0，Ace 点数为 1。
func (c Card) PointValue() int {
	if c.Rank >= Ten {
		return 0
	}
	return int(c.Rank)
}

// String returns a short string representation of the card (e.g., "AS", "10H").
// String 返回该牌的简短字符串表示（例如："AS"、"10H"）。
func (c Card) String() string {
	ranks := map[Rank]string{
		Ace: "A", Two: "2", Three: "3", Four: "4", Five: "5", Six: "6",
		Seven: "7", Eight: "8", Nine: "9", Ten: "10", Jack: "J", Queen: "Q", King: "K",
	}
	suits := map[Suit]string{
		Spades: "♠", Hearts: "♥", Diamonds: "♦", Clubs: "♣",
	}
	return ranks[c.Rank] + suits[c.Suit]
}
