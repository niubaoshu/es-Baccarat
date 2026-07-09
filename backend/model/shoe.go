package model

import (
	"errors"
	"math/rand"
	"time"
)

var ErrShoeEmpty = errors.New("shoe is empty")
var ErrPastCutCard = errors.New("cut card reached, please shuffle shoe")

// Shoe represents the dealer's shoe containing multiple decks of cards.
// Shoe 代表发牌靴，内含多副牌。
type Shoe struct {
	Cards            []Card //牌堆
	DecksCount       int    //牌堆数
	CutCardThreshold int    //切牌阈值
	currentIndex     int    //当前索引
}

// NewShoe initializes a new Shoe with a basic, unshuffled set of decks.
// NewShoe 初始化一个新的发牌靴，包含未洗牌的基础牌组。
func NewShoe(decksCount int, cutCardThreshold int) *Shoe {
	s := &Shoe{
		DecksCount:       decksCount,
		CutCardThreshold: cutCardThreshold,
		Cards:            make([]Card, 0, decksCount*52),
	}
	s.populate()
	return s
}

// populate fills the shoe with standard decks in order.
// populate 按顺序将标准牌组填充到发牌靴中。
func (s *Shoe) populate() {
	s.Cards = s.Cards[:0]
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	for d := 0; d < s.DecksCount; d++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				s.Cards = append(s.Cards, Card{Suit: suit, Rank: rank})
			}
		}
	}
	s.currentIndex = 0
}

// Shuffle randomizes the order of the cards in the shoe and resets the current index.
// Shuffle 将发牌靴中的牌随机打乱，并重置当前索引。
func (s *Shoe) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(s.Cards) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i]
	}
	s.currentIndex = 0
}

// Draw returns the next card from the shoe. Returns ErrShoeEmpty if there are no cards left.
// Note: It does NOT return an error specifically if past the cut card, it merely allows drawing.
// Use IsPastCutCard() to check if a new shoe should be prepared for the *next* round.
// Draw 从发牌靴中取出下一张牌。如果牌已取完，则返回 ErrShoeEmpty。
// 注意：即使超过切牌位置，该函数也不会返回错误，仍允许继续取牌。
// 请使用 IsPastCutCard() 检查是否需要为*下一局*准备新的发牌靴。
func (s *Shoe) Draw() (Card, error) {
	if s.currentIndex >= len(s.Cards) {
		return Card{}, ErrShoeEmpty
	}
	c := s.Cards[s.currentIndex]
	s.currentIndex++
	return c, nil
}

// CardsLeft returns the number of cards remaining in the shoe.
// CardsLeft 返回发牌靴中剩余的牌数。
func (s *Shoe) CardsLeft() int {
	return len(s.Cards) - s.currentIndex
}

// IsPastCutCard returns true if the number of cards left is less than or equal to the CutCardThreshold.
// In actual gameplay, if this returns true the current hand is finished, and a new shoe/shuffle is triggered before the next hand.
// IsPastCutCard 如果剩余牌数小于或等于切牌阈值，则返回 true。
// 在实际游戏中，若返回 true，则当前局结束后，在下一局开始前需重新洗牌或更换发牌靴。
func (s *Shoe) IsPastCutCard() bool {
	return s.CardsLeft() <= s.CutCardThreshold
}

// Burn performs the standard Baccarat burn card procedure.
// It draws one face up card, looks at its baccarat point value (10/J/Q/K is considered 10 for burning purposes in many casinos),
// and then burns (draws and discards) that many cards.
// Burn 执行百家乐标准的烧牌流程。
// 先翻开一张牌，查看其百家乐点数（在许多赌场中，10/J/Q/K 的烧牌点数视为 10），
// 然后烧掉（取出并弃置）相应数量的牌。
func (s *Shoe) Burn() error {
	faceUpCard, err := s.Draw()
	if err != nil {
		return err
	}

	burnCount := int(faceUpCard.Rank)
	if faceUpCard.Rank >= Ten {
		burnCount = 10
	}

	for i := 0; i < burnCount; i++ {
		_, err := s.Draw()
		if err != nil {
			return err
		}
	}
	return nil
}
