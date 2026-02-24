package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/niubaoshu/es-Baccarat/backend/config"
	"github.com/niubaoshu/es-Baccarat/backend/model"
	"github.com/niubaoshu/es-Baccarat/backend/rules"
)

// SimulationStats holds the aggregated results of a simulation run.
type SimulationStats struct {
	TotalRounds  int
	OutcomeCount map[rules.Outcome]int
	Duration     time.Duration
}

// RunSimulation executes a fast, headless Monte Carlo simulation of Baccarat.
func RunSimulation(cfg *config.GameConfig, totalRounds int, numWorkers int) *SimulationStats {
	start := time.Now()

	// Adjust workers if needed
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if totalRounds < numWorkers {
		numWorkers = totalRounds
	}

	roundsPerWorker := totalRounds / numWorkers
	remainder := totalRounds % numWorkers

	var wg sync.WaitGroup
	resultsCh := make(chan map[rules.Outcome]int, numWorkers)

	fmt.Printf("Starting simulation of %d rounds using %d workers...\n", totalRounds, numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		targetRounds := roundsPerWorker
		if w == 0 {
			targetRounds += remainder // First worker takes the remainder
		}

		go func(rounds int) {
			defer wg.Done()

			localCounts := make(map[rules.Outcome]int)
			shoe := model.NewShoe(cfg.DecksCount, cfg.CutCardThreshold)
			shoe.Shuffle()
			_ = shoe.Burn()

			for i := 0; i < rounds; i++ {
				// Re-shoe if needed
				if shoe.IsPastCutCard() {
					shoe = model.NewShoe(cfg.DecksCount, cfg.CutCardThreshold)
					shoe.Shuffle()
					_ = shoe.Burn()
				}

				// Deal
				c1, _ := shoe.Draw()
				c2, _ := shoe.Draw()
				c3, _ := shoe.Draw()
				c4, _ := shoe.Draw()

				pHand := &model.Hand{Cards: []model.Card{c1, c3}}
				bHand := &model.Hand{Cards: []model.Card{c2, c4}}

				var pThird *model.Card
				playerHit := rules.DeterminePlayerHit(pHand, bHand)
				if playerHit {
					c, _ := shoe.Draw()
					pHand.AddCard(c)
					pThird = &c
				}

				bankerHit := rules.DetermineBankerHit(bHand, pHand, playerHit, pThird)
				if bankerHit {
					c, _ := shoe.Draw()
					bHand.AddCard(c)
				}

				outcome := rules.DetermineOutcome(pHand, bHand)
				localCounts[outcome]++
			}

			resultsCh <- localCounts
		}(targetRounds)
	}

	wg.Wait()
	close(resultsCh)

	finalCounts := make(map[rules.Outcome]int)
	for res := range resultsCh {
		for outcome, count := range res {
			finalCounts[outcome] += count
		}
	}

	return &SimulationStats{
		TotalRounds:  totalRounds,
		OutcomeCount: finalCounts,
		Duration:     time.Since(start),
	}
}

// PrintReport prints the statistical percentages to the console.
func (s *SimulationStats) PrintReport() {
	fmt.Printf("\n=== Simulation Complete ===\n")
	fmt.Printf("Total Rounds: %d\n", s.TotalRounds)
	fmt.Printf("Time Taken:   %s (%.0f rounds/sec)\n", s.Duration, float64(s.TotalRounds)/s.Duration.Seconds())
	fmt.Println()
	totalPlayerWins := s.OutcomeCount[rules.OutcomePlayer] + s.OutcomeCount[rules.OutcomePanda8]
	totalCount := totalPlayerWins + s.OutcomeCount[rules.OutcomeBanker] + s.OutcomeCount[rules.OutcomeTie] + s.OutcomeCount[rules.OutcomeDragon7]

	pPlayerSim := (float64(totalPlayerWins) / float64(s.TotalRounds)) * 100
	pPandaSim := (float64(s.OutcomeCount[rules.OutcomePanda8]) / float64(s.TotalRounds)) * 100
	pBankerSim := (float64(s.OutcomeCount[rules.OutcomeBanker]) / float64(s.TotalRounds)) * 100
	pTieSim := (float64(s.OutcomeCount[rules.OutcomeTie]) / float64(s.TotalRounds)) * 100
	pDragonSim := (float64(s.OutcomeCount[rules.OutcomeDragon7]) / float64(s.TotalRounds)) * 100
	pTotalSim := (float64(totalCount) / float64(s.TotalRounds)) * 100

	fmt.Printf("%-20s | %-12s | %-12s | %-12s\n", "Outcome", "Count", "Simulated %", "Expected %")
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "Player (Total)", totalPlayerWins, pPlayerSim, 44.6247)
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "  ↳ Panda 8", s.OutcomeCount[rules.OutcomePanda8], pPandaSim, 3.4543)
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "Banker (Non-Dragon)", s.OutcomeCount[rules.OutcomeBanker], pBankerSim, 43.6064)
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "Tie", s.OutcomeCount[rules.OutcomeTie], pTieSim, 9.5156)
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "Dragon 7", s.OutcomeCount[rules.OutcomeDragon7], pDragonSim, 2.2534)
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%-20s | %12d | %11.4f%% | %11.4f%%\n", "Total", totalCount, pTotalSim, 100.0000)
	fmt.Printf("==================================================================\n")

	// Calculate simulated EV for $1 bets on each option
	betProfits := map[rules.BetType]int{
		rules.Player: 0,
		rules.Banker: 0,
		rules.Tie:    0,
		rules.Dragon: 0,
		rules.Panda:  0,
	}

	for outcome, count := range s.OutcomeCount {
		for bet := range betProfits {
			payout := rules.CalculatePayout(outcome, bet, 1)
			netChange := payout.NetChange(1)
			betProfits[bet] += netChange * count
		}
	}

	fmt.Printf("\n%-20s | %-16s | %-15s | %-15s\n", "Bet Type ($1/hand)", "Net Profit ($)", "Simulated EV", "Expected EV")
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Printf("%-20s | %16d | %14.4f%% | %14.4f%%\n", "Banker", betProfits[rules.Banker], float64(betProfits[rules.Banker])/float64(s.TotalRounds)*100, -1.0183)
	fmt.Printf("%-20s | %16d | %14.4f%% | %14.4f%%\n", "Player", betProfits[rules.Player], float64(betProfits[rules.Player])/float64(s.TotalRounds)*100, -1.2351)
	fmt.Printf("%-20s | %16d | %14.4f%% | %14.4f%%\n", "Tie", betProfits[rules.Tie], float64(betProfits[rules.Tie])/float64(s.TotalRounds)*100, -14.3596)
	fmt.Printf("%-20s | %16d | %14.4f%% | %14.4f%%\n", "Dragon 7", betProfits[rules.Dragon], float64(betProfits[rules.Dragon])/float64(s.TotalRounds)*100, -7.6106)
	fmt.Printf("%-20s | %16d | %14.4f%% | %14.4f%%\n", "Panda 8", betProfits[rules.Panda], float64(betProfits[rules.Panda])/float64(s.TotalRounds)*100, -10.1882)
	fmt.Printf("=======================================================================\n\n")
}
