package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/niubaoshu/es-Baccarat/config"
	"github.com/niubaoshu/es-Baccarat/engine"
	"github.com/niubaoshu/es-Baccarat/player"
)

func main() {
	var (
		playerName      string
		createPlayer    bool
		initialBalance  int
		simulateRounds  int
		simulateWorkers int
	)

	flag.StringVar(&playerName, "player", "", "Specify the player username")
	flag.BoolVar(&createPlayer, "create_player", false, "Create a new player profile")
	flag.IntVar(&initialBalance, "initial_balance", 10000, "Initial balance for a new player (default 10000)")
	flag.IntVar(&simulateRounds, "simulate", 0, "Number of rounds to simulate mathematically (if > 0, skips interactive mode)")
	flag.IntVar(&simulateWorkers, "workers", 4, "Number of concurrent workers for simulation")
	flag.Parse()

	cfg := config.DefaultConfig()

	// --- Simulation Mode ---
	if simulateRounds > 0 {
		stats := engine.RunSimulation(cfg, simulateRounds, simulateWorkers)
		stats.PrintReport()
		return
	}

	// --- Interactive Mode ---
	if playerName == "" && !createPlayer {
		playerName = "default_player"
		fmt.Printf("No player specified. Using '%s'.\n", playerName)

		// Attempt to create implicitly if doesn't exist
		_, err := player.LoadProfile(playerName)
		if err == player.ErrPlayerNotFound {
			_, err = player.CreateProfile(playerName, initialBalance)
			if err != nil {
				fmt.Printf("Fatal error creating default player: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Implicitly created '%s' with balance %d.\n", playerName, initialBalance)
		}
	} else if playerName == "" && createPlayer {
		fmt.Println("Error: --player must be provided when using --create_player")
		flag.Usage()
		os.Exit(1)
	}

	// 2. Profile Loading or Creation
	var p *player.Profile
	var err error

	if createPlayer {
		p, err = player.CreateProfile(playerName, initialBalance)
		if err != nil {
			if err == player.ErrPlayerAlreadyExists {
				fmt.Printf("Error: Player '%s' already exists. Cannot recreate or overwrite balance.\n", playerName)
			} else {
				fmt.Printf("Error creating player: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Successfully created '%s' with starting balance %d.\n", playerName, initialBalance)
	} else {
		p, err = player.LoadProfile(playerName)
		if err != nil {
			if err == player.ErrPlayerNotFound {
				fmt.Printf("Error: Player '%s' not found. Please use --create_player to register.\n", playerName)
			} else {
				fmt.Printf("Error loading profile: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Welcome back, %s! Loaded historical balance: $%d (Total Hands: %d)\n", p.Username, p.Balance, p.HandsPlayed)
	}

	// 3. Initialize Game Engine
	game := engine.NewGame(cfg, p)

	// 4. Main Game Loop
	fmt.Println("\n--- Starting EZ Baccarat Session ---")
	for {
		fmt.Printf("\n[ Current Balance: $%d ]\n", game.Profile.Balance)
		if game.Profile.Balance <= 0 {
			fmt.Println("You are out of money! Game Over.")
			break
		}

		bets := engine.PromptBets()
		if bets == nil {
			fmt.Println("Thanks for playing! Exiting...")
			break
		}

		totalBetAmount := 0
		for _, v := range bets {
			totalBetAmount += v
		}

		if totalBetAmount > game.Profile.Balance {
			fmt.Printf("Error: Insufficient funds. Total bet ($%d) exceeds balance ($%d).\n", totalBetAmount, game.Profile.Balance)
			continue
		}

		game.PlayRound(bets)
	}
}
