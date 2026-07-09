package config

// GameConfig holds the core settings for the Baccarat simulator.
// GameConfig 保存百家乐模拟器的核心配置项。
type GameConfig struct {
	DecksCount       int // 牌堆数
	CutCardThreshold int // 切牌阈值
}

// DefaultConfig returns the standard casino settings.
// DefaultConfig 返回标准赌场的配置参数。
func DefaultConfig() *GameConfig {
	return &GameConfig{
		DecksCount:       8,
		CutCardThreshold: 14, // Roughly 1/4 of a deck
		// 大约为一副牌的 1/4
	}
}
