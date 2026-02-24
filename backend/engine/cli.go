package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/niubaoshu/es-Baccarat/backend/rules"
)

// PromptBets asks the user to enter their bets via the terminal.
func PromptBets() map[rules.BetType]int {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your bets for this round.")
	fmt.Println("Available types: P (Player), B (Banker), T (Tie), D (Dragon 7), 8 (Panda 8).")
	fmt.Println("Format: <Type>:<Amount> separate multiple by comma. (e.g. P:100,D:10)")
	fmt.Println("Leave empty to stop playing (Quit).")

	for {
		fmt.Print("Your bets: ")
		if !scanner.Scan() {
			return nil
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" || strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" {
			return nil
		}

		parts := strings.Split(input, ",")
		parsedBets := make(map[rules.BetType]int)
		valid := true

		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			kv := strings.Split(p, ":")
			if len(kv) != 2 {
				fmt.Println("Invalid format. Use <Type>:<Amount>")
				valid = false
				break
			}

			bTypeStr := strings.ToUpper(strings.TrimSpace(kv[0]))
			amtStr := strings.TrimSpace(kv[1])

			amt, err := strconv.Atoi(amtStr)
			if err != nil || amt <= 0 {
				fmt.Printf("Invalid amount: %s\n", amtStr)
				valid = false
				break
			}

			var bType rules.BetType
			switch bTypeStr {
			case "P", "PLAYER":
				bType = rules.Player
			case "B", "BANKER":
				bType = rules.Banker
			case "T", "TIE":
				bType = rules.Tie
			case "D", "DRAGON", "DRAGON7", "DRAGON 7":
				bType = rules.Dragon
			case "8", "PANDA", "PANDA8", "PANDA 8":
				bType = rules.Panda
			default:
				fmt.Printf("Unknown bet type: %s\n", bTypeStr)
				valid = false
			}

			if valid {
				parsedBets[bType] += amt
			}
		}

		if valid && len(parsedBets) > 0 {
			// Rule validation: Panda and Dragon requires Player or Banker bet.
			// Same for Tie if strict, but the new document says Tie can be placed independently.
			hasBase := parsedBets[rules.Player] > 0 || parsedBets[rules.Banker] > 0
			hasSpecial := parsedBets[rules.Dragon] > 0 || parsedBets[rules.Panda] > 0

			if hasSpecial && !hasBase {
				fmt.Println("Rule Error: Dragon 7 and Panda 8 bets require an active Player or Banker base bet.")
				continue
			}

			return parsedBets
		}
	}
}
