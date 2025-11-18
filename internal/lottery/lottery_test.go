package lottery

import (
	"fmt"
	"testing"

	"github.com/palemoky/lucky-day/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateWeight(t *testing.T) {
	// 使用 testify 的 InDelta 来处理浮点数比较，更加健壮
	const delta = 0.001 // 定义一个可接受的误差范围

	testCases := []struct {
		name           string
		participant    model.Participant
		currentYear    int
		expectedWeight float64
	}{
		{
			name:           "没有中奖记录的新用户",
			participant:    model.Participant{ID: 1, Name: "Newbie"},
			currentYear:    2025,
			expectedWeight: 1.0,
		},
		{
			name: "去年中过三等奖的用户 (权重触底)",
			participant: model.Participant{
				ID: 2, Name: "LastYearWinner",
				WinningHistory: []model.WinningRecord{{Year: 2024, PrizeLevel: 3}},
			},
			currentYear:    2025,
			expectedWeight: 0.01,
		},
		{
			name: "五年前中过特等奖的用户",
			participant: model.Participant{
				ID: 3, Name: "OldWinner",
				WinningHistory: []model.WinningRecord{{Year: 2020, PrizeLevel: 0}},
			},
			currentYear:    2025,
			expectedWeight: 0.3767,
		},
		{
			name: "去年中过特等奖的用户 (权重触底)",
			participant: model.Participant{
				ID: 4, Name: "HeavyPenalty",
				WinningHistory: []model.WinningRecord{{Year: 2024, PrizeLevel: 0}},
			},
			currentYear:    2025,
			expectedWeight: 0.01,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weight := calculateWeight(tc.participant, tc.currentYear)
			// 使用 InDelta 断言实际值在期望值的正负 delta 范围内
			assert.InDelta(t, tc.expectedWeight, weight, delta, "计算出的权重应在期望值附近")
		})
	}
}

// helper function to create a standard set of participants for tests
func createTestParticipants(count int) []model.Participant {
	participants := make([]model.Participant, count)
	for i := range count {
		participants[i] = model.Participant{ID: i + 1, Name: fmt.Sprintf("User%d", i+1)}
	}
	return participants
}

// helper function to create a standard set of prizes for tests
func createTestPrizes() []model.Prize {
	return []model.Prize{
		{ID: 1, Name: "一等奖", Count: 1},
		{ID: 2, Name: "二等奖", Count: 3},
		{ID: 3, Name: "三等奖", Count: 10},
	}
}

func TestEngine_Draw(t *testing.T) {
	testCases := []struct {
		name              string
		participantsCount int
		prizeIDToDraw     int
		expectedWinCount  int
		expectedEligible  int
		expectSuccess     bool
	}{
		{
			name:              "正常抽取二等奖",
			participantsCount: 20,
			prizeIDToDraw:     2, // 二等奖，3个名额
			expectedWinCount:  3,
			expectedEligible:  17, // 20 - 3
			expectSuccess:     true,
		},
		{
			name:              "候选人不足时全部中奖",
			participantsCount: 2,
			prizeIDToDraw:     2, // 二等奖，3个名额
			expectedWinCount:  2,
			expectedEligible:  0, // 2 - 2
			expectSuccess:     true,
		},
		{
			name:              "抽取不存在的奖项",
			participantsCount: 20,
			prizeIDToDraw:     99, // 不存在的奖项ID
			expectedWinCount:  0,
			expectedEligible:  20,
			expectSuccess:     false,
		},
		{
			name:              "连续抽取",
			participantsCount: 20,
			prizeIDToDraw:     1, // 先抽一等奖
			expectedWinCount:  1,
			expectedEligible:  19,
			expectSuccess:     true,
		},
		{
			name:              "候选人池为空时抽奖",
			participantsCount: 0,
			prizeIDToDraw:     1,
			expectedWinCount:  0,
			expectedEligible:  0,
			expectSuccess:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 每个测试用例都使用全新的Engine实例
			participants := createTestParticipants(tc.participantsCount)
			prizes := createTestPrizes()
			engine := NewEngine(participants, prizes)

			winners, ok := engine.Draw(tc.prizeIDToDraw)

			// 使用 testify/assert 进行断言
			assert.Equal(t, tc.expectSuccess, ok)
			assert.Len(t, winners, tc.expectedWinCount, "中奖人数应该符合预期")
			assert.Len(t, engine.eligible, tc.expectedEligible, "剩余候选人数应该符合预期")

			if ok {
				// 验证中奖者确实从候选池中移除了
				for _, winner := range winners {
					_, exists := engine.eligible[winner.ID]
					assert.False(t, exists, "中奖者不应再存在于候选池中")
				}
			}

			// 特殊处理“连续抽取”用例
			if tc.name == "连续抽取" {
				// 接着抽二等奖
				winners2, ok2 := engine.Draw(2)
				assert.True(t, ok2)
				assert.Len(t, winners2, 3)
				assert.Len(t, engine.eligible, 16) // 19 - 3
			}
		})
	}
}

func TestEngine_ResetPrize(t *testing.T) {
	// 1. Setup: 创建一个有10个参与者的引擎
	participants := createTestParticipants(10)
	prizes := createTestPrizes()
	engine := NewEngine(participants, prizes)

	// 2. Action: 抽取二等奖 (3名)
	winners, ok := engine.Draw(2)
	require.True(t, ok) // require 会在失败时立即停止测试，适合前置条件
	require.Len(t, winners, 3)
	require.Len(t, engine.eligible, 7)
	require.Equal(t, 3, engine.prizes[1].DrawnCount) // 假设prizes[1]是二等奖

	// 3. Action: 重置二等奖
	engine.ResetPrize(2)

	// 4. Assert: 验证状态是否已恢复
	assert.Len(t, engine.eligible, 10, "重置后，所有参与者都应有资格")

	// 验证 prize 的 DrawnCount 是否已重置
	var prize2_reset model.Prize
	for _, p := range engine.GetPrizes() {
		if p.ID == 2 {
			prize2_reset = p
			break
		}
	}
	assert.Equal(t, 0, prize2_reset.DrawnCount, "奖品的已抽取数量应重置为0")

	// 验证 allWinners map 中已无该奖项记录
	_, ok = engine.allWinners[2]
	assert.False(t, ok, "总中奖名单中不应再有该奖项的记录")

	// 5. 验证可以重新抽取该奖项
	winners_after_reset, ok_after_reset := engine.Draw(2)
	assert.True(t, ok_after_reset)
	assert.Len(t, winners_after_reset, 3)
	assert.Len(t, engine.eligible, 7)
}

func FuzzEngine_Draw(f *testing.F) {
	// 1. 添加种子语料库 (seed corpus)
	f.Add(100, 5, 10) // 100人，5个奖品，每个奖品10个名额
	f.Add(10, 20, 1)  // 人比奖品名额少
	f.Add(1000, 1, 100)
	f.Add(5, 2, 3)

	// 2. 编写模糊测试的目标函数
	f.Fuzz(func(t *testing.T, numParticipants int, numPrizes int, prizeCount int) {
		// 3. 数据校验和准备
		// Fuzzing 会产生随机输入，包括负数或0，需要保护测试逻辑本身不 panic
		if numParticipants <= 0 || numPrizes <= 0 || prizeCount <= 0 {
			t.Skip() // 跳过无效的测试输入
		}
		if numParticipants > 10000 || numPrizes > 100 || prizeCount > 100 {
			t.Skip() // 避免过大的输入导致测试超时
		}

		participants := make([]model.Participant, numParticipants)
		for i := range numParticipants {
			participants[i] = model.Participant{ID: i, Name: fmt.Sprintf("FuzzUser%d", i)}
		}

		prizes := make([]model.Prize, numPrizes)
		for i := range numPrizes {
			prizes[i] = model.Prize{ID: i, Name: fmt.Sprintf("FuzzPrize%d", i), Count: prizeCount}
		}

		engine := NewEngine(participants, prizes)

		// 4. 执行并断言
		// 循环抽取所有奖品
		currentEligibleCount := numParticipants
		for i := range numPrizes {
			// 如果已经没有候选人了，就没必要继续抽了
			if currentEligibleCount == 0 {
				break
			}

			// 计算本轮期望抽出的中奖人数
			// 应该是 奖品名额 和 剩余候选人数 中的较小值
			expectedDrawCount := prizeCount
			if currentEligibleCount < expectedDrawCount {
				expectedDrawCount = currentEligibleCount
			}

			winners, ok := engine.Draw(prizes[i].ID)

			// 断言本次抽奖的结果
			assert.True(t, ok, "当有候选人时，抽奖应该返回 true")
			assert.Len(t, winners, expectedDrawCount, "抽出中奖者的人数应符合预期")

			// 更新我们追踪的剩余人数
			currentEligibleCount -= len(winners)

			// 验证引擎的内部状态是否与我们的追踪一致
			assert.Len(t, engine.GetEligibleParticipants(), currentEligibleCount, "最终剩余候选人数应正确")
		}
	})
}
