package card

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
type card struct {
	suit Suit
	rank Rank
}

// PointValue calculates the Baccarat point value of the card.
// 2-9 retain their value. 10, J, Q, K are worth 0. Ace is worth 1.
// PointValue 计算该牌在百家乐中的点数值。
// 2-9 保留其原始点数，10、J、Q、K 点数为 0，Ace 点数为 1。
func (c card) pointValue() int {
	if c.rank >= Ten {
		return 0
	}
	return int(c.rank)
}

// String returns a short string representation of the card (e.g., "AS", "10H").
// String 返回该牌的简短字符串表示（例如："AS"、"10H"）。
func (c card) String() string {
	ranks := map[Rank]string{
		Ace: "A", Two: "2", Three: "3", Four: "4", Five: "5", Six: "6",
		Seven: "7", Eight: "8", Nine: "9", Ten: "10", Jack: "J", Queen: "Q", King: "K",
	}
	suits := map[Suit]string{
		Spades: "♠", Hearts: "♥", Diamonds: "♦", Clubs: "♣",
	}
	return ranks[c.rank] + suits[c.suit]
}

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

// Hand represents a player's or banker's hand of cards in Baccarat.
// Hand 代表百家乐中闲家或庄家的一手牌。
type Hand struct {
	cards []card
}

func NewHand() *Hand {
	return &Hand{
		cards: make([]card, 0, 3),
	}
}

func (h *Hand) Reset() {
	h.cards = h.cards[:0]
}

// addCard adds a card to the hand.
// AddCard 向手牌中添加一张牌。
func (h *Hand) addCard(c card) {
	h.cards = append(h.cards, c)
}

// TotalPoints returns the Baccarat point value of the hand (0-9).
// It sums the point values of all cards and drops the tens digit.
// TotalPoints 返回该手牌的百家乐点数（0-9）。
// 它将所有牌的点数相加，并取个位数。
func (h *Hand) totalPoints() int {
	total := 0
	for _, c := range h.cards {
		total += c.pointValue()
	}
	return total % 10
}

// isNatural checks if the hand is a "Natural" win (total of 8 or 9 in the first two cards).
// isNatural 检查该手牌是否为"天牌"（前两张牌合计为 8 或 9）。
func (h *Hand) isNatural() bool {
	if len(h.cards) == 2 {
		points := h.totalPoints()
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
	for i, c := range h.cards {
		if i > 0 {
			res += ", "
		}
		res += c.String()
	}
	return res
}

// Round represent one game
type Round struct {
	bHand *Hand
	pHand *Hand
}

func (r *Round) Reset() {
	r.bHand.Reset()
	r.pHand.Reset()
}

// determinePlayerHit decides if the Player should draw a third card.
// According to Baccarat rules, if either has a Natural (8 or 9), no one hits.
// If not natural, Player hits on 0-5, stands on 6-9.
// determinePlayerHit 决定闲家是否应该补第三张牌。
// 根据百家乐规则，若任意一方为天牌（8 或 9），则双方均不补牌。
// 若非天牌，闲家点数 0-5 时补牌，6-9 时不补。

// 闲家补牌规则
// 首轮发牌后，由闲家先执行规则。判断依据仅为闲家前两张牌的总点数和。闲家前两张牌点数和执行动作
// 8 或 9（天牌/例牌）双方均直接停牌，比大小判定胜负
// 6 或 7闲家停牌（不补牌）
// 0、1、2、3、4、5闲家必须补发第三张牌
func (r *Round) determinePlayerHit() bool {
	if r.pHand.isNatural() || r.bHand.isNatural() {
		return false
	}

	if r.pHand.totalPoints() <= 5 {
		return true
	}
	return false
}

// determineBankerHit decides if the Banker should draw a third card.
// It requires knowing whether the Player has already hit, and what specific card they drew.
// If either has a Natural (8 or 9), no one hits.
// determineBankerHit 决定庄家是否应该补第三张牌。
// 需要知道闲家是否已补牌，以及闲家所补的具体牌面。
// 若任意一方为天牌（8 或 9），则双方均不补牌。
func (r *Round) determineBankerHit() bool {
	if r.pHand.isNatural() || r.bHand.isNatural() {
		return false
	}

	bankerPts := r.bHand.totalPoints()
	playerHit := len(r.pHand.cards) == 3
	var playerThirdCard *card

	// If Player DID hit, the Banker drawing depends on the third card drawn by the player
	// 若闲家已补牌，庄家是否补牌取决于闲家所补的第三张牌
	if playerHit {
		playerThirdCard = &r.pHand.cards[2]

		p3Pts := playerThirdCard.pointValue()

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
		case 4:
			// Banker hits on 4 if player's third card is 2-7
			// 庄家点数为 4 且闲家第三张牌为 2-7 时补牌
			if p3Pts >= 2 && p3Pts <= 7 {
				return true
			}
		case 5:
			// Banker hits on 5 if player's third card is 4-7
			// 庄家点数为 5 且闲家第三张牌为 4-7 时补牌
			if p3Pts >= 4 && p3Pts <= 7 {
				return true
			}
		case 6:
			// Banker hits on 6 if player's third card is 6 or 7
			// 庄家点数为 6 且闲家第三张牌为 6 或 7 时补牌
			if p3Pts == 6 || p3Pts == 7 {
				return true
			}
		}
	} else {
		// If player did NOT hit (stood on 6 or 7)
		// 若闲家未补牌（停在 6 或 7）

		// Banker hits on 0-5, stands on 6-7
		// 庄家点数 0-5 时补牌，6-7 时不补
		return bankerPts <= 5
	}
	return false
}

// DetermineOutcome compares final hands and returns the specific EZ Baccarat outcome.
// DetermineOutcome 比较最终手牌并返回具体的 EZ 百家乐结果。
func (r *Round) DetermineOutcome() Outcome {
	pPts := r.pHand.totalPoints()
	bPts := r.bHand.totalPoints()

	if pPts == bPts {
		return OutcomeTie
	}

	if pPts > bPts {
		// Player Wins. Check for Panda 8
		// 闲家赢。检查是否为熊猫8
		if pPts == 8 && len(r.pHand.cards) == 3 {
			return OutcomePanda8
		}
		return OutcomePlayer
	}

	// Banker Wins. Check for Dragon 7
	// 庄家赢。检查是否为龙7
	if bPts == 7 && len(r.bHand.cards) == 3 {
		return OutcomeDragon7
	}
	return OutcomeBanker
}
