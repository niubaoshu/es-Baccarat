package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// RoundLog defines what gets written to the logging file for every hand played.
type RoundLog struct {
	Timestamp      time.Time      `json:"timestamp"`
	Player         string         `json:"player"`
	InitialBalance int            `json:"initial_balance"`
	FinalBalance   int            `json:"final_balance"`
	Bets           map[string]int `json:"bets"` // BetType -> Amount
	PlayerHand     []string       `json:"player_hand"`
	BankerHand     []string       `json:"banker_hand"`
	PlayerPoints   int            `json:"player_points"`
	BankerPoints   int            `json:"banker_points"`
	Outcome        string         `json:"outcome"`
	NetChange      int            `json:"net_change"`
}

var logDir = "data/logs"

// LogRound appends a round summary to the JSONL log file.
func LogRound(logEntry RoundLog) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "game_history.jsonl")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}
