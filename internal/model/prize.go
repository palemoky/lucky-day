package model

// PrizeLevel 定义奖品等级类型
type PrizeLevel int

const (
	PrizeLevelSpecial PrizeLevel = iota // 0: 特等奖
	PrizeLevel1                         // 1: 一等奖
	PrizeLevel2                         // 2: 二等奖
	PrizeLevel3                         // 3: 三等奖
	PrizeLevel4                         // 4: 四等奖
	PrizeLevel5                         // 5: 五等奖
)

func (pl PrizeLevel) String() string {
	switch pl {
	case PrizeLevelSpecial:
		return "特等奖"
	case PrizeLevel1:
		return "一等奖"
	case PrizeLevel2:
		return "二等奖"
	case PrizeLevel3:
		return "三等奖"
	case PrizeLevel4:
		return "四等奖"
	case PrizeLevel5:
		return "五等奖"
	default:
		return "未知奖项"
	}
}

// Prize 奖品结构体
type Prize struct {
	ID          int
	Name        string
	Level       PrizeLevel
	Count       int     // 奖品数量
	Probability float64 // 中奖概率
	DrawnCount  int     // 已抽奖数量
}

// Participant 参与者结构体
type Participant struct {
	ID             int `gorm:"primaryKey"`
	Name           string
	WinningHistory []WinningRecord `gorm:"foreignKey:ParticipantID"`
}

// WinningRecord 往年中奖记录
type WinningRecord struct {
	ID            uint `gorm:"primaryKey"`
	ParticipantID int
	Year          int
	PrizeLevel    int
}
