package config

// GameConfig holds the core settings for the Baccarat simulator.
type GameConfig struct {
	DecksCount       int
	CutCardThreshold int
}

// DefaultConfig returns the standard casino settings.
func DefaultConfig() *GameConfig {
	return &GameConfig{
		DecksCount:       8,
		CutCardThreshold: 14, // Roughly 1/4 of a deck
	}
}
