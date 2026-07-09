package card

import (
	"math/rand"
	"time"
)

// shoe represents the dealer's shoe containing multiple decks of cards.
// shoe 代表发牌靴，内含多副牌。
type Shoe struct {
	cards            []card //牌堆
	decksCount       int    //牌堆数
	cutCardThreshold int    //切牌阈值
	currentIndex     int    //当前索引
	seed             int64  //随机数种子
	round            *Round //上一局游戏
}

// NewShoe initializes a new shoe with a basic, unshuffled set of decks.
// seed 用于洗牌的随机数种子，传入 0 则使用当前时间作为种子。
// NewShoe 初始化一个新的发牌靴，包含未洗牌的基础牌组。
func NewShoe(decksCount int, cutCardThreshold int, seed int64) *Shoe {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	s := &Shoe{
		decksCount:       decksCount,
		cutCardThreshold: cutCardThreshold,
		cards:            make([]card, 0, decksCount*52),
		seed:             seed,
	}
	s.populate()
	s.shuffle()
	s.burn()

	pHand := NewHand()
	bHand := NewHand()
	s.round = &Round{
		pHand: pHand,
		bHand: bHand,
	}

	return s
}

func (s *Shoe) Reset() {
	s.shuffle()
	s.burn()
}

func (s *Shoe) NextRound() *Round {
	s.round.Reset()
	s.round.pHand.addCard(s.draw())
	s.round.bHand.addCard(s.draw())
	s.round.pHand.addCard(s.draw())
	s.round.bHand.addCard(s.draw())

	if s.round.determinePlayerHit() {
		s.round.pHand.addCard(s.draw())
	}

	if s.round.determineBankerHit() {
		s.round.bHand.addCard(s.draw())
	}

	return s.round
}

// populate fills the shoe with standard decks in order.
// populate 按顺序将标准牌组填充到发牌靴中。
func (s *Shoe) populate() {
	s.cards = s.cards[:0]
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	for d := 0; d < s.decksCount; d++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				s.cards = append(s.cards, card{suit: suit, rank: rank})
			}
		}
	}
	s.currentIndex = 0
}

// shuffle randomizes the order of the cards in the shoe and resets the current index.
// shuffle 将发牌靴中的牌随机打乱，并重置当前索引。
func (s *Shoe) shuffle() {
	r := rand.New(rand.NewSource(s.seed))
	for i := len(s.cards) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	}
	s.currentIndex = 0
}

// draw returns the next card from the shoe.
// draw 从发牌靴中取出下一张牌。
func (s *Shoe) draw() card {
	c := s.cards[s.currentIndex]
	s.currentIndex++
	return c
}

// burn performs the standard Baccarat burn card procedure.
// It draws one face up card, looks at its baccarat point value (10/J/Q/K is considered 10 for burning purposes in many casinos),
// and then burns (draws and discards) that many cards.
// burn 执行百家乐标准的烧牌流程。
// 先翻开一张牌，查看其百家乐点数（在许多赌场中，10/J/Q/K 的烧牌点数视为 10），
// 然后烧掉（取出并弃置）相应数量的牌。
func (s *Shoe) burn() {
	burnCount := int(s.draw().rank)
	if burnCount >= 10 {
		burnCount = 10
	}

	for i := 0; i < burnCount; i++ {
		s.draw()
	}
}

// CardsLeft returns the number of cards remaining in the shoe.
// CardsLeft 返回发牌靴中剩余的牌数。
func (s *Shoe) cardsLeft() int {
	return len(s.cards) - s.currentIndex
}

// IsPastCutCard returns true if the number of cards left is less than or equal to the CutCardThreshold.
// In actual gameplay, if this returns true the current hand is finished, and a new shoe/shuffle is triggered before the next hand.
// IsPastCutCard 如果剩余牌数小于或等于切牌阈值，则返回 true。
// 在实际游戏中，若返回 true，则当前局结束后，在下一局开始前需重新洗牌或更换发牌靴。

// 在发牌过程中抽出割牌：如果荷官在发当前局的牌时抽出了割牌，这一局会继续正常发完并结算。
// 因为割牌后面还有足够的备用牌。此局结算后，当前牌靴结束，更换新靴。

// 在发牌前露出割牌：如果上一局刚结束，牌靴最前面露出来的第一张正好是割牌，荷官会宣布“最后一局”（Last Hand）不成立，
// 直接结束这靴牌，不再进行下一局的投注和发牌。
func (s *Shoe) IsPastCutCard() bool {
	return s.cardsLeft() <= s.cutCardThreshold
}
