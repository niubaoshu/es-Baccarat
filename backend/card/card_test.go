package card

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestCardExample(t *testing.T) {
	// 创建一个 8 副牌、切牌阈值 14 的牌靴；seed=0 表示使用当前时间作为随机种子
	shoe := NewShoe(8, 14, 0)

	fmt.Printf("牌靴初始剩余牌数: %d\n\n", shoe.cardsLeft())

	// 真实游戏循环：每局前检查切牌线，过线则换靴
	round := 0
	for !shoe.IsPastCutCard() {
		round++
		r := shoe.NextRound()
		outcome := r.DetermineOutcome()

		fmt.Printf("第 %02d 局 | 闲家: %-16s %d 点 | 庄家: %-16s %d 点 | 结果: %-10s | 剩余: %d 张\n",
			round,
			r.pHand.String(), r.pHand.totalPoints(),
			r.bHand.String(), r.bHand.totalPoints(),
			outcome,
			shoe.cardsLeft(),
		)
	}

	fmt.Printf("\n共发 %d 局，牌靴到达切牌线，剩余 %d 张，需换靴。\n", round, shoe.cardsLeft())
}

// newManualShoe 直接用指定牌序构造 shoe，绕过洗牌和烧牌，用于规则验证。
// newManualShoe constructs a shoe directly from the given card sequence,
// bypassing shuffle and burn, for rule verification purposes.
func newManualShoe(cards []card) *Shoe {
	return &Shoe{
		cards:            cards,
		decksCount:       1,
		cutCardThreshold: 0,
		currentIndex:     0,
		round:            &Round{pHand: NewHand(), bHand: NewHand()},
	}
}

// TestCardRuleVerification 手动构造牌组，逻辑推理 EZ 百家乐规则，验证代码行为。
// 发牌顺序: cards[0]=闲p1, [1]=庄b1, [2]=闲p2, [3]=庄b2, [4]=第5张, [5]=第6张
func TestCardRuleVerification(t *testing.T) {
	cases := []struct {
		name     string
		cards    []card  // 按发牌顺序排列
		expected Outcome // 逻辑推理出的预期结果
		reason   string  // 推理过程说明
	}{
		{
			name: "闲家天牌9",
			// 闲家: 5+4=9（天牌）→ 双方不补牌
			// 庄家: 3+2=5
			// 闲9 > 庄5，闲家2张，不触发 Panda8 → OutcomePlayer
			cards:    []card{{Spades, Five}, {Diamonds, Three}, {Hearts, Four}, {Clubs, Two}},
			expected: OutcomePlayer,
			reason:   "闲家天牌9(5+4)，庄家5点，闲胜，非熊猫8（只有2张）",
		},
		{
			name: "庄家天牌8",
			// 闲家: 2+3=5
			// 庄家: 4+4=8（天牌）→ 双方不补牌
			// 庄8 > 闲5，庄家2张，不触发 Dragon7 → OutcomeBanker
			cards:    []card{{Spades, Two}, {Diamonds, Four}, {Hearts, Three}, {Clubs, Four}},
			expected: OutcomeBanker,
			reason:   "庄家天牌8(4+4)，闲家5点，庄胜，非龙7（只有2张）",
		},
		{
			name: "龙7：庄家三张合计7点获胜",
			// 闲家: 2+4=6 → 停牌（6点停）
			// 庄家: 1+2=3 → 闲家未补，庄家3<=5 → 补牌
			// 庄家补牌: 4 → 1+2+4=7，三张
			// 庄7 > 闲6，庄家3张且7点 → OutcomeDragon7
			cards:    []card{{Spades, Two}, {Diamonds, Ace}, {Hearts, Four}, {Clubs, Two}, {Spades, Four}},
			expected: OutcomeDragon7,
			reason:   "闲家6点停牌，庄家3点补牌后得7(1+2+4)三张，龙7",
		},
		{
			name: "熊猫8：闲家三张合计8点获胜",
			// 闲家: 3+2=5 → 补牌
			// 庄家: 3+3=6
			// 闲家补牌: 3 → 3+2+3=8，三张
			// 庄家: 闲家第三张=3点，庄6点时只在第三张6或7时补牌 → 不补
			// 闲8 > 庄6，闲家3张且8点 → OutcomePanda8
			cards:    []card{{Spades, Three}, {Diamonds, Three}, {Hearts, Two}, {Clubs, Three}, {Hearts, Three}},
			expected: OutcomePanda8,
			reason:   "闲家5点补牌后得8(3+2+3)三张，庄家6点第三张非6/7不补，熊猫8",
		},
		{
			name: "平局",
			// 闲家: 3+4=7 → 停牌
			// 庄家: 2+5=7 → 闲家未补，庄家7>5 → 停牌
			// 闲7 == 庄7 → OutcomeTie
			cards:    []card{{Spades, Three}, {Diamonds, Two}, {Hearts, Four}, {Clubs, Five}},
			expected: OutcomeTie,
			reason:   "闲家7点停牌，庄家7点停牌，平局",
		},
		{
			name: "闲庄各补一张，庄家胜（非龙7）",
			// 闲家: 2+3=5 → 补牌
			// 庄家: 1+3=4
			// 闲家补牌: 5 → 2+3+5=10→0点
			// 庄家: 闲第三张=5点，庄4点时在2-7补牌 → 补牌
			// 庄家补牌: 1 → 1+3+1=5点
			// 闲0 < 庄5，庄家3张但5点，非龙7 → OutcomeBanker
			cards:    []card{{Spades, Two}, {Diamonds, Ace}, {Hearts, Three}, {Clubs, Three}, {Spades, Five}, {Clubs, Ace}},
			expected: OutcomeBanker,
			reason:   "闲家补5后得0点，庄家补A后得5点，庄胜，庄家3张但5点非龙7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newManualShoe(tc.cards)
			r := s.NextRound()
			got := r.DetermineOutcome()

			// 打印详情方便对照
			fmt.Printf("\n[%s]\n  推理: %s\n  闲家: %s (%d点, %d张)\n  庄家: %s (%d点, %d张)\n  预期: %-12s  实际: %s\n",
				tc.name, tc.reason,
				r.pHand.String(), r.pHand.totalPoints(), len(r.pHand.cards),
				r.bHand.String(), r.bHand.totalPoints(), len(r.bHand.cards),
				tc.expected, got,
			)

			if got != tc.expected {
				t.Errorf("结果不符: 预期 %q，实际 %q", tc.expected, got)
			}
		})
	}
}

// TestMultiRoundShoe 用一个手工构造的长牌序 shoe，连续发 10 局。
// 每局先用文字推理 EZ 百家乐规则得出预期，再与代码结果对照。
//
// 发牌顺序（每局）：p1, b1, p2, b2, [p3（若闲补）], [b3（若庄补）]
// 补牌规则速查：
//
//	闲家：天牌(8/9)→停，6/7→停，0-5→补
//	庄家（闲未补）：0-5→补，6/7→停
//	庄家（闲已补，看闲第3张点数）：
//	  庄0-2→始终补  庄3→非8补  庄4→2-7补  庄5→4-7补  庄6→6-7补  庄7→停
//	特殊牌型：闲3张8点→Panda8，庄3张7点→Dragon7
func TestMultiRoundShoe(t *testing.T) {
	// ═══════════════════════════════════════════════════════════════════
	// 牌序设计（10 局，共 50 张）
	// ═══════════════════════════════════════════════════════════════════
	//
	// 局1 (4张): 闲天牌9 → Player
	//   p1=6♠ b1=2♦ p2=3♥ b2=4♣
	//   闲: 6+3=9 → 天牌，双方不补
	//   庄: 2+4=6
	//   闲9 > 庄6，2张，非Panda8 → Player
	//
	// 局2 (4张): 庄天牌8 → Banker
	//   p1=2♠ b1=5♦ p2=3♥ b2=3♣
	//   闲: 2+3=5
	//   庄: 5+3=8 → 天牌，双方不补
	//   庄8 > 闲5，2张，非Dragon7 → Banker
	//
	// 局3 (4张): 平局，双方停牌 → Tie
	//   p1=3♠ b1=4♦ p2=4♥ b2=3♣
	//   闲: 3+4=7 → 停牌
	//   庄: 4+3=7 → 闲未补，庄7>5 → 停牌
	//   闲7 == 庄7 → Tie
	//
	// 局4 (5张): 闲停，庄补，庄赢（非Dragon7） → Banker
	//   p1=3♠ b1=A♦ p2=3♥ b2=2♣ b3=3♦
	//   闲: 3+3=6 → 停牌（6点停）
	//   庄: 1+2=3 → 闲未补，庄3<=5 → 补牌
	//   庄补3♦: 1+2+3=6点，3张
	//   庄6 == 闲6 → Tie ← 注意！平局，不是庄赢
	//   修正：改庄补牌让庄>闲
	//   b3改为4♦: 庄1+2+4=7点，庄7>闲6，庄3张7点 → Dragon7
	//
	// 局4 (5张): 闲6停，庄3点补4→7 → Dragon7
	//   p1=3♠ b1=A♦ p2=3♥ b2=2♣ b3=4♦
	//   闲: 3+3=6 → 停牌
	//   庄: 1+2=3 → 补; 1+2+4=7，3张
	//   庄7 > 闲6，庄3张7点 → Dragon7
	//
	// 局5 (5张): 闲补，Panda8 → Panda8
	//   p1=2♠ b1=3♦ p2=3♥ b2=3♣ p3=3♠
	//   闲: 2+3=5 → 补牌
	//   庄: 3+3=6
	//   闲补3♠: 2+3+3=8，3张
	//   庄: 闲第3张=3点，庄6点只在6/7时补 → 不补
	//   闲8 > 庄6，闲3张8点 → Panda8
	//
	// 局6 (6张): 闲补，庄也补，闲赢（非Panda8） → Player
	//   p1=2♠ b1=2♦ p2=3♥ b2=2♣ p3=2♦ b3=2♠
	//   闲: 2+3=5 → 补牌
	//   庄: 2+2=4
	//   闲补2♦: 2+3+2=7点，3张（非8点，非Panda8）
	//   庄: 闲第3张=2点，庄4点在2-7时补 → 补牌
	//   庄补2♠: 2+2+2=6点
	//   闲7 > 庄6，闲3张7点（非Panda8） → Player
	//
	// 局7 (4张): 闲天牌8 → Player
	//   p1=5♠ b1=2♦ p2=3♥ b2=4♣
	//   闲: 5+3=8 → 天牌，双方不补
	//   庄: 2+4=6
	//   闲8 > 庄6，2张，非Panda8 → Player
	//
	// 局8 (6张): 闲补庄补，庄赢（3张，非Dragon7） → Banker
	//   p1=2♠ b1=A♦ p2=3♥ b2=3♣ p3=5♠ b3=A♥
	//   闲: 2+3=5 → 补牌
	//   庄: 1+3=4
	//   闲补5♠: 2+3+5=10→0点
	//   庄: 闲第3张=5点，庄4点在2-7时补 → 补牌
	//   庄补A♥: 1+3+1=5点，3张
	//   庄5 > 闲0，庄3张5点（非Dragon7） → Banker
	//
	// 局9 (5张): 闲补，庄不补，庄赢 → Banker
	//   p1=3♠ b1=4♦ p2=2♥ b2=4♣ p3=6♠
	//   闲: 3+2=5 → 补牌
	//   庄: 4+4=8 → 天牌！双方不补
	//   → 闲天牌检查：闲是否天牌？3+2=5，非天牌
	//   → 庄天牌8：isNatural=true → determinePlayerHit返回false → 闲不补
	//   修正：庄有天牌→闲不补→playerHit=false
	//   重新：闲3+2=5（非天牌），庄4+4=8（天牌）→ 闲不补，庄不补
	//   庄8 > 闲5 → Banker（2张，非Dragon7）
	//   这局实际只用4张牌（p3=6♠不会被取到）
	//   注意：p3=6♠放在牌序里但不会被抽取，它变成下一局的p1！
	//
	//   修正局9：确保闲不是天牌、庄是天牌，只消耗4张
	//   p1=2♠ b1=5♦ p2=3♥ b2=3♣
	//   闲:2+3=5，庄:5+3=8(天牌)→闲不补、庄不补，庄8>闲5 → Banker
	//   → 但这和局2一样了。换个花色即可（不影响逻辑，只是牌不同）
	//
	// 局10 (4张): 闲7点停，庄4点停（闲未补时庄4点不足5），庄补后...
	//   实际：闲未补，庄4点→ 4<=5 → 庄补牌
	//   重新设计局10：闲7停，庄6停（>5不补）→ 闲7>庄6 → Player
	//   p1=4♠ b1=3♦ p2=3♥ b2=3♣
	//   闲: 4+3=7 → 停牌
	//   庄: 3+3=6 → 闲未补，庄6>5 → 停牌
	//   闲7 > 庄6 → Player
	// ═══════════════════════════════════════════════════════════════════

	shoecards := []card{
		// 局1 (4张): 闲天牌9 → Player
		{Spades, Six}, {Diamonds, Two}, {Hearts, Three}, {Clubs, Four},
		// 局2 (4张): 庄天牌8 → Banker
		{Spades, Two}, {Diamonds, Five}, {Hearts, Three}, {Clubs, Three},
		// 局3 (4张): 双方停，平局 → Tie
		{Spades, Three}, {Diamonds, Four}, {Hearts, Four}, {Clubs, Three},
		// 局4 (5张): 闲6停，庄3点补4→7(3张) → Dragon7
		{Spades, Three}, {Diamonds, Ace}, {Hearts, Three}, {Clubs, Two}, {Diamonds, Four},
		// 局5 (5张): 闲5补3→8(3张) Panda8
		{Spades, Two}, {Diamonds, Three}, {Hearts, Three}, {Clubs, Three}, {Spades, Three},
		// 局6 (6张): 闲5补2→7(3张)，庄4补2→6，闲7>庄6 → Player
		{Spades, Two}, {Diamonds, Two}, {Hearts, Three}, {Clubs, Two}, {Diamonds, Two}, {Spades, Two},
		// 局7 (4张): 闲天牌8 → Player
		{Spades, Five}, {Diamonds, Two}, {Hearts, Three}, {Clubs, Four},
		// 局8 (6张): 闲5补5→0，庄4补A→5(3张) → Banker
		{Spades, Two}, {Diamonds, Ace}, {Hearts, Three}, {Clubs, Three}, {Spades, Five}, {Hearts, Ace},
		// 局9 (4张): 庄天牌8，闲5→庄胜 → Banker
		{Spades, Two}, {Diamonds, Five}, {Hearts, Three}, {Clubs, Three},
		// 局10 (4张): 闲7停，庄6停 → Player
		{Spades, Four}, {Diamonds, Three}, {Hearts, Three}, {Clubs, Three},
	}

	// 推理汇总表
	expected := []struct {
		name   string
		reason string
		want   Outcome
	}{
		{
			"局1",
			"闲6+3=9天牌，庄2+4=6，双方不补，闲9>庄6，2张非Panda8",
			OutcomePlayer,
		},
		{
			"局2",
			"闲2+3=5，庄5+3=8天牌，双方不补，庄8>闲5，2张非Dragon7",
			OutcomeBanker,
		},
		{
			"局3",
			"闲3+4=7停，庄4+3=7停（闲未补，庄7>5不补），7==7平局",
			OutcomeTie,
		},
		{
			"局4",
			"闲3+3=6停，庄1+2=3补4→7(3张)，庄7>闲6，Dragon7",
			OutcomeDragon7,
		},
		{
			"局5",
			"闲2+3=5补3→8(3张)，庄3+3=6(第3张=3不在6/7)不补，闲8>庄6，Panda8",
			OutcomePanda8,
		},
		{
			"局6",
			"闲2+3=5补2→7(3张)，庄2+2=4(第3张=2在2-7)补2→6，闲7>庄6，3张非Panda8",
			OutcomePlayer,
		},
		{
			"局7",
			"闲5+3=8天牌，庄2+4=6，双方不补，闲8>庄6，2张非Panda8",
			OutcomePlayer,
		},
		{
			"局8",
			"闲2+3=5补5→0，庄1+3=4(第3张=5在2-7)补A→5(3张)，庄5>闲0，3张非Dragon7",
			OutcomeBanker,
		},
		{
			"局9",
			"闲2+3=5，庄5+3=8天牌，双方不补，庄8>闲5，2张非Dragon7",
			OutcomeBanker,
		},
		{
			"局10",
			"闲4+3=7停，庄3+3=6停（闲未补，庄6>5不补），闲7>庄6",
			OutcomePlayer,
		},
	}

	s := newManualShoe(shoecards)
	fmt.Printf("\n%-6s %-50s %-12s %-12s %s\n", "局数", "推理", "预期", "实际", "状态")
	fmt.Println("────────────────────────────────────────────────────────────────────────────────────────────")

	for i, exp := range expected {
		r := s.NextRound()
		got := r.DetermineOutcome()
		status := "✅ PASS"
		if got != exp.want {
			status = "❌ FAIL"
		}

		fmt.Printf("%-6s %-50s %-12s %-12s %s\n", exp.name, exp.reason, exp.want, got, status)
		fmt.Printf("       闲: %-14s %d点 %d张   庄: %-14s %d点 %d张\n",
			r.pHand.String(), r.pHand.totalPoints(), len(r.pHand.cards),
			r.bHand.String(), r.bHand.totalPoints(), len(r.bHand.cards),
		)

		if got != exp.want {
			t.Errorf("第%d局结果不符: 预期 %q, 实际 %q", i+1, exp.want, got)
		}
	}
}

// TestSimulation100M 模拟 1 亿局 EZ 百家乐，统计各种结果的出现概率。
// 运行时请加大超时：go test ./card/... -v -run TestSimulation100M -count=1 -timeout 300s
func TestSimulation100M(t *testing.T) {
	const totalRounds = 100_000_000

	var (
		countPlayer  int64
		countBanker  int64
		countTie     int64
		countDragon7 int64
		countPanda8  int64
	)

	// 用一个独立的 rng 为每次换靴生成新种子，保证每靴洗牌独立
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	s := NewShoe(8, 14, rng.Int63())

	start := time.Now()

	for i := 0; i < totalRounds; i++ {
		if !s.IsPastCutCard() {
			r := s.NextRound()
			switch r.DetermineOutcome() {
			case OutcomePlayer:
				countPlayer++
			case OutcomeBanker:
				countBanker++
			case OutcomeTie:
				countTie++
			case OutcomeDragon7:
				countDragon7++
			case OutcomePanda8:
				countPanda8++
			}
		} else {
			s.Reset()
		}
	}

	elapsed := time.Since(start)
	total := int64(totalRounds)

	fmt.Printf("\n模拟 %d 局完成，耗时 %v\n", totalRounds, elapsed.Round(time.Millisecond))
	fmt.Printf("平均每局耗时: %.0f ns\n\n", float64(elapsed.Nanoseconds())/float64(total))
	fmt.Printf("%-14s %12s %10s\n", "结果", "局数", "概率")
	fmt.Println("───────────────────────────────────────")

	printRow := func(name string, count int64) {
		fmt.Printf("%-14s %12d %9.5f%%\n", name, count, float64(count)/float64(total)*100)
	}
	printRow("闲家胜 (Player)", countPlayer)
	printRow("庄家胜 (Banker)", countBanker)
	printRow("平局   (Tie)", countTie)
	printRow("龙7 (Dragon 7)", countDragon7)
	printRow("熊猫8 (Panda 8)", countPanda8)
	fmt.Println("───────────────────────────────────────")

	// 闲胜合计 = Player + Panda8，庄胜合计 = Banker + Dragon7
	printRow("闲胜合计", countPlayer+countPanda8)
	printRow("庄胜合计", countBanker+countDragon7)
	printRow("合计", countPlayer+countBanker+countTie+countDragon7+countPanda8)
}
