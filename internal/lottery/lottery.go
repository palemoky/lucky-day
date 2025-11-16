package lottery

import (
	"math"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand"
	"github.com/palemoky/lucky-day/internal/model"
)

// Engine 抽奖引擎，封装了状态和逻辑
type Engine struct {
	allParticipants []model.Participant
	prizes          []model.Prize
	eligible        map[int]model.Participant   // 仍有资格抽奖的参与者
	allWinners      map[int][]model.Participant // 所有奖项的中奖者，Key 是 Prize.ID
}

func NewEngine(participants []model.Participant, prizes []model.Prize) *Engine {
	eligibleMap := make(map[int]model.Participant)
	for _, p := range participants {
		eligibleMap[p.ID] = p
	}

	return &Engine{
		allParticipants: participants,
		prizes:          prizes,
		eligible:        eligibleMap,
		allWinners:      make(map[int][]model.Participant),
	}
}

// Draw 为指定奖项抽出中奖者
func (e *Engine) Draw(prizeID int) ([]model.Participant, bool) {
	var prizeToDraw *model.Prize
	prizeIndex := -1
	for i, p := range e.prizes {
		if p.ID == prizeID {
			prizeToDraw = &e.prizes[i]
			prizeIndex = i
			break
		}
	}
	if prizeToDraw == nil {
		return nil, false // 奖项不存在
	}

	// 如果该奖项名额已满，则不允许再抽
	if prizeToDraw.DrawnCount >= prizeToDraw.Count {
		return nil, false
	}

	// 确定本次需要抽取的人数
	drawCount := prizeToDraw.Count - prizeToDraw.DrawnCount

	// 构造权重选择器
	choices := e.getWeightedChoices(*prizeToDraw)
	if len(choices) == 0 {
		return nil, false // 没有可抽奖的人了
	}
	// 如果候选人数少于等于要抽取的人数，则全部中奖
	if len(choices) <= drawCount {
		winners := make([]model.Participant, 0, len(choices))
		for _, choice := range choices {
			winners = append(winners, choice.Item.(model.Participant))
		}
		// 更新状态
		for _, winner := range winners {
			delete(e.eligible, winner.ID)
		}
		e.allWinners[prizeToDraw.ID] = append(e.allWinners[prizeToDraw.ID], winners...)
		e.prizes[prizeIndex].DrawnCount += len(winners)

		return winners, true
	}

	// 如果候选人多于要抽取的人数，则开始抽奖
	chooser, _ := weightedrand.NewChooser(choices...)
	winnersMap := make(map[int]model.Participant) // 使用 map 来确保中奖者不重复

	// 循环抽奖，直到抽满 drawCount 个不重复的中奖者
	for len(winnersMap) < drawCount {
		winner := chooser.Pick().(model.Participant)
		winnersMap[winner.ID] = winner
	}

	var currentWinners []model.Participant
	for _, winner := range winnersMap {
		currentWinners = append(currentWinners, winner)
		// 从总候选池中移除中奖者
		delete(e.eligible, winner.ID)
	}

	// 更新中奖记录和奖品已抽取数量
	e.allWinners[prizeToDraw.ID] = append(e.allWinners[prizeToDraw.ID], currentWinners...)
	e.prizes[prizeIndex].DrawnCount += len(currentWinners)

	return currentWinners, true
}

// GetEligibleParticipants 获取当前所有有资格的参与者
func (e *Engine) GetEligibleParticipants() []model.Participant {
	var participants []model.Participant
	for _, p := range e.eligible {
		participants = append(participants, p)
	}
	return participants
}

// getWeightedChoices 为当前有资格的参与者生成加权选项
func (e *Engine) getWeightedChoices(prize model.Prize) []weightedrand.Choice {
	currentYear := time.Now().Year()
	var choices []weightedrand.Choice

	// 如果候选人池为空，直接返回
	if len(e.eligible) == 0 {
		return choices
	}

	for _, participant := range e.eligible {
		// 乘以1000以提高权重计算的精度
		weight := uint(calculateWeight(participant, currentYear) * prize.Probability * 1000)
		if weight > 0 {
			choices = append(choices, weightedrand.Choice{Item: participant, Weight: weight})
		}
	}
	// 如果计算后所有人的权重都是0，则给予每个人相同的权重
	if len(choices) == 0 {
		for _, participant := range e.eligible {
			choices = append(choices, weightedrand.Choice{Item: participant, Weight: 100})
		}
	}

	return choices
}

// calculateWeight 计算参与者的最终抽奖权重 (此函数逻辑来自你的版本)
func calculateWeight(participant model.Participant, currentYear int) float64 {
	baseWeight := 1.0
	decayFactor := 0.5        // 衰减因子
	levelPenaltyFactor := 1.5 // 等级惩罚因子

	for _, record := range participant.WinningHistory {
		yearDiff := float64(currentYear - record.Year)
		timePenalty := math.Exp(-decayFactor * yearDiff)                           // 时间惩罚
		levelPenalty := math.Pow(levelPenaltyFactor, float64(5-record.PrizeLevel)) // 等级惩罚
		baseWeight -= timePenalty * levelPenalty
	}

	if baseWeight < 0.01 {
		return 0.01 // 保证最低权重
	}
	return baseWeight
}

// GetPrizes 返回奖品列表
func (e *Engine) GetPrizes() []model.Prize {
	return e.prizes
}

// GetAllWinners 返回所有中奖者
func (e *Engine) GetAllWinners() map[int][]model.Participant {
	return e.allWinners
}

// ResetPrize 重置某个奖项，使其可以重新抽取
func (e *Engine) ResetPrize(prizeID int) {
	// 1. 将该奖项的中奖者放回 eligible 池
	if winners, ok := e.allWinners[prizeID]; ok {
		for _, winner := range winners {
			e.eligible[winner.ID] = winner
		}
	}

	// 2. 清空该奖项的中奖记录
	delete(e.allWinners, prizeID)

	// 3. 重置奖品的 DrawnCount
	for i := range e.prizes {
		if e.prizes[i].ID == prizeID {
			e.prizes[i].DrawnCount = 0
			break
		}
	}
}

// GetRandomNames 从有资格的参与者中随机挑选N个名字用于动画
func (e *Engine) GetRandomNames(count int) []string {
	eligible := e.GetEligibleParticipants()
	if len(eligible) == 0 {
		return []string{"无候选人"}
	}

	names := make([]string, count)
	for i := range count {
		randIndex := rand.Intn(len(eligible))
		names[i] = eligible[randIndex].Name
	}
	return names
}
