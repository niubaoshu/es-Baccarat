package engine

import (
	"fmt"
	"time"

	"github.com/niubaoshu/es-Baccarat/config"
	"github.com/niubaoshu/es-Baccarat/model"
	"github.com/niubaoshu/es-Baccarat/player"
	"github.com/niubaoshu/es-Baccarat/rules"
)

// Game orchestrates the physical simulation of the shoe and hands.
type Game struct {
	Config  *config.GameConfig
	Shoe    *model.Shoe
	Profile *player.Profile
}

// NewGame initializes a game session.
func NewGame(cfg *config.GameConfig, p *player.Profile) *Game {
	g := &Game{
		Config:  cfg,
		Profile: p,
	}
	g.initShoe()
	return g
}

func (g *Game) initShoe() {
	fmt.Printf("\n[Dealer] Bringing out a new shoe with %d decks...\n", g.Config.DecksCount)
	g.Shoe = model.NewShoe(g.Config.DecksCount, g.Config.CutCardThreshold)
	g.Shoe.Shuffle()
	fmt.Println("[Dealer] Shuffling cards...")

	err := g.Shoe.Burn()
	if err != nil {
		fmt.Printf("[Error] Failed to burn cards: %v\n", err)
	} else {
		fmt.Println("[Dealer] Burn procedure complete.")
	}
}

// PlayRound handles the end-to-end logic for a single round of Baccarat given user bets.
func (g *Game) PlayRound(bets map[rules.BetType]int) {
	if g.Shoe.IsPastCutCard() {
		fmt.Println("\n[Dealer] Cut card reached. Preparing new shoe...")
		g.initShoe()
	}

	initialBalance := g.Profile.Balance
	totalBetAmount := 0
	for _, amt := range bets {
		totalBetAmount += amt
	}

	// 1. Deduct bets
	g.Profile.Balance -= totalBetAmount
	g.Profile.TotalWager += totalBetAmount
	g.Profile.HandsPlayed++

	// 2. Deal initial cards
	pHand := &model.Hand{}
	bHand := &model.Hand{}

	c1, _ := g.Shoe.Draw() // Player 1
	c2, _ := g.Shoe.Draw() // Banker 1
	c3, _ := g.Shoe.Draw() // Player 2
	c4, _ := g.Shoe.Draw() // Banker 2

	pHand.AddCard(c1)
	bHand.AddCard(c2)
	pHand.AddCard(c3)
	bHand.AddCard(c4)

	fmt.Printf("\n--- [Deal Completed] ---\n")
	fmt.Printf("Player Hand: %s  (Total: %d)\n", pHand.String(), pHand.TotalPoints())
	fmt.Printf("Banker Hand: %s  (Total: %d)\n", bHand.String(), bHand.TotalPoints())

	// 3. Process Third Card Rules
	var pThirdCard *model.Card
	playerHit := rules.DeterminePlayerHit(pHand, bHand)

	if playerHit {
		c, _ := g.Shoe.Draw()
		fmt.Printf("[Action] Player hits and draws: %s\n", c.String())
		pHand.AddCard(c)
		pThirdCard = &c
		fmt.Printf("Player Final Hand: %s  (Total: %d)\n", pHand.String(), pHand.TotalPoints())
	} else if pHand.IsNatural() || bHand.IsNatural() {
		fmt.Println("[Action] Natural 8 or 9 detected. No hits.")
	} else {
		fmt.Println("[Action] Player stands.")
	}

	bankerHit := rules.DetermineBankerHit(bHand, pHand, playerHit, pThirdCard)

	if bankerHit {
		c, _ := g.Shoe.Draw()
		fmt.Printf("[Action] Banker hits and draws: %s\n", c.String())
		bHand.AddCard(c)
		fmt.Printf("Banker Final Hand: %s  (Total: %d)\n", bHand.String(), bHand.TotalPoints())
	} else if !pHand.IsNatural() && !bHand.IsNatural() {
		fmt.Println("[Action] Banker stands.")
	}

	// 4. Outcomes and Payouts
	outcome := rules.DetermineOutcome(pHand, bHand)
	fmt.Printf("\n>>> [Outcome]: %s Wins! <<<\n", outcome)

	totalWin := 0
	totalReturned := 0

	for bType, amt := range bets {
		result := rules.CalculatePayout(outcome, bType, amt)
		totalWin += result.WinAmount
		totalReturned += result.Returned

		change := result.NetChange(amt)
		if change > 0 {
			fmt.Printf("  - %s Bet ($%d): WIN (+%d)\n", bType, amt, result.WinAmount)
		} else if change == 0 {
			fmt.Printf("  - %s Bet ($%d): PUSH\n", bType, amt)
		} else {
			fmt.Printf("  - %s Bet ($%d): LOSE\n", bType, amt)
		}
	}

	g.Profile.Balance += totalWin + totalReturned
	netChange := (totalWin + totalReturned) - totalBetAmount

	// 5. Save State and Log
	_ = g.Profile.Save()

	pStringCards := make([]string, len(pHand.Cards))
	for i, c := range pHand.Cards {
		pStringCards[i] = c.String()
	}

	bStringCards := make([]string, len(bHand.Cards))
	for i, c := range bHand.Cards {
		bStringCards[i] = c.String()
	}

	strBets := make(map[string]int)
	for k, v := range bets {
		strBets[string(k)] = v
	}

	log := RoundLog{
		Timestamp:      time.Now(),
		Player:         g.Profile.Username,
		InitialBalance: initialBalance,
		FinalBalance:   g.Profile.Balance,
		Bets:           strBets,
		PlayerHand:     pStringCards,
		BankerHand:     bStringCards,
		PlayerPoints:   pHand.TotalPoints(),
		BankerPoints:   bHand.TotalPoints(),
		Outcome:        string(outcome),
		NetChange:      netChange,
	}
	_ = LogRound(log)

	// 6. Round Summary Print
	fmt.Printf("\n=== Round Summary ===\n")
	fmt.Printf("Cards Left: %d\n", g.Shoe.CardsLeft())
	fmt.Printf("Net Change: $%d\n", netChange)
	fmt.Printf("New Balance: $%d\n", g.Profile.Balance)
	fmt.Printf("=====================\n\n")
}
