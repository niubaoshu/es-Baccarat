package model

import (
	"errors"
	"math/rand"
	"time"
)

var ErrShoeEmpty = errors.New("shoe is empty")
var ErrPastCutCard = errors.New("cut card reached, please shuffle shoe")

// Shoe represents the dealer's shoe containing multiple decks of cards.
type Shoe struct {
	Cards            []Card
	DecksCount       int
	CutCardThreshold int
	currentIndex     int
}

// NewShoe initializes a new Shoe with a basic, unshuffled set of decks.
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
func (s *Shoe) Draw() (Card, error) {
	if s.currentIndex >= len(s.Cards) {
		return Card{}, ErrShoeEmpty
	}
	c := s.Cards[s.currentIndex]
	s.currentIndex++
	return c, nil
}

// CardsLeft returns the number of cards remaining in the shoe.
func (s *Shoe) CardsLeft() int {
	return len(s.Cards) - s.currentIndex
}

// IsPastCutCard returns true if the number of cards left is less than or equal to the CutCardThreshold.
// In actual gameplay, if this returns true the current hand is finished, and a new shoe/shuffle is triggered before the next hand.
func (s *Shoe) IsPastCutCard() bool {
	return s.CardsLeft() <= s.CutCardThreshold
}

// Burn performs the standard Baccarat burn card procedure.
// It draws one face up card, looks at its baccarat point value (10/J/Q/K is considered 10 for burning purposes in many casinos),
// and then burns (draws and discards) that many cards.
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
