package datasource

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/model"
)

// LoadParticipants 统一入口函数
func LoadParticipants(cfg config.DataSourceConfig) ([]model.Participant, error) {
	switch cfg.Type {
	case "csv":
		fmt.Println("数据源: CSV 文件")
		return loadParticipantsFromCSV(cfg.CSV.Path)
	case "excel":
		fmt.Println("数据源: Excel 文件")
		return LoadParticipantsFromExcel(cfg.Excel.Path)
	case "db":
		fmt.Printf("数据源: 数据库 (%s)\n", cfg.Database.Driver)
		return loadParticipantsFromDB(cfg.Database)
	default:
		return nil, fmt.Errorf("未知的数据源类型: %s", cfg.Type)
	}
}

// loadParticipantsFromCSV
func loadParticipantsFromCSV(filePath string) ([]model.Participant, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 CSV 文件: %w", err)
	}
	defer func() { _ = file.Close() }() // Ignore close error

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("无法读取 CSV 内容: %w", err)
	}

	if len(records) <= 1 {
		return []model.Participant{}, nil
	}

	var participants []model.Participant
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}
		id, _ := strconv.Atoi(record[0])
		participant := model.Participant{
			ID:             id,
			Name:           record[1],
			WinningHistory: []model.WinningRecord{}, // CSV 简化处理
		}
		participants = append(participants, participant)
	}
	fmt.Printf("成功从 CSV 文件 [%s] 加载了 %d 名参与者。\n", filePath, len(participants))
	return participants, nil
}
