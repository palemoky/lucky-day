package datasource

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/palemoky/lucky-day/internal/model"
)

const (
	SheetPrizes       = "Prizes"
	SheetParticipants = "Participants"
	SheetWinners      = "Winners"
)

// LoadPrizesFromExcel loads prizes from Excel file
func LoadPrizesFromExcel(filePath string) ([]model.Prize, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close Excel file: %v\n", err)
		}
	}()

	rows, err := f.GetRows(SheetPrizes)
	if err != nil {
		return nil, fmt.Errorf("failed to read Prizes sheet: %w", err)
	}

	if len(rows) <= 1 {
		return nil, fmt.Errorf("prizes sheet is empty or only contains header")
	}

	var prizes []model.Prize
	for i, row := range rows[1:] { // Skip header
		if len(row) < 6 {
			continue // Skip incomplete rows
		}

		var id, count, level int
		var probability float64

		// Parse ID
		if _, err := fmt.Sscanf(row[0], "%d", &id); err != nil {
			fmt.Printf("warning: skipping row %d, invalid ID: %v\n", i+2, err)
			continue
		}

		// Parse Count
		if _, err := fmt.Sscanf(row[3], "%d", &count); err != nil {
			fmt.Printf("warning: skipping row %d, invalid Count: %v\n", i+2, err)
			continue
		}

		// Parse Level
		if _, err := fmt.Sscanf(row[4], "%d", &level); err != nil {
			fmt.Printf("warning: skipping row %d, invalid Level: %v\n", i+2, err)
			continue
		}

		// Parse Probability
		if _, err := fmt.Sscanf(row[5], "%f", &probability); err != nil {
			fmt.Printf("warning: skipping row %d, invalid Probability: %v\n", i+2, err)
			continue
		}

		prize := model.Prize{
			ID:          id,
			Name:        row[1], // Use Chinese name by default
			Level:       model.PrizeLevel(level),
			Count:       count,
			Probability: probability,
			DrawnCount:  0,
		}
		prizes = append(prizes, prize)
	}

	fmt.Printf("Successfully loaded %d prizes from Excel file [%s]\n", len(prizes), filePath)
	return prizes, nil
}

// LoadParticipantsFromExcel loads participants from Excel file
func LoadParticipantsFromExcel(filePath string) ([]model.Participant, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close Excel file: %v\n", err)
		}
	}()

	rows, err := f.GetRows(SheetParticipants)
	if err != nil {
		return nil, fmt.Errorf("failed to read Participants sheet: %w", err)
	}

	if len(rows) <= 1 {
		return nil, fmt.Errorf("participants sheet is empty or only contains header")
	}

	var participants []model.Participant
	for i, row := range rows[1:] { // Skip header
		if len(row) < 2 {
			continue // Skip incomplete rows
		}

		var id int
		if _, err := fmt.Sscanf(row[0], "%d", &id); err != nil {
			fmt.Printf("warning: skipping row %d, invalid ID: %v\n", i+2, err)
			continue
		}

		participant := model.Participant{
			ID:             id,
			Name:           row[1],
			WinningHistory: []model.WinningRecord{}, // Will be loaded separately if needed
		}
		participants = append(participants, participant)
	}

	fmt.Printf("Successfully loaded %d participants from Excel file [%s]\n", len(participants), filePath)
	return participants, nil
}

// Winner represents a lottery winner for Excel export
type Winner struct {
	DrawTime   time.Time
	PrizeName  string
	WinnerID   int
	WinnerName string
	PrizeLevel int
}

// SaveWinnersToExcel saves winners to Excel file
func SaveWinnersToExcel(filePath string, winners []Winner) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close Excel file: %v\n", err)
		}
	}()

	// Get existing rows to append new winners
	rows, err := f.GetRows(SheetWinners)
	if err != nil {
		return fmt.Errorf("failed to read Winners sheet: %w", err)
	}

	startRow := len(rows) + 1
	if len(rows) == 0 {
		// If sheet is empty, add header first
		header := []interface{}{"Draw Time", "Prize Name", "Winner ID", "Winner Name", "Prize Level"}
		if err := f.SetSheetRow(SheetWinners, "A1", &header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		startRow = 2
	}

	// Write winners
	for i, winner := range winners {
		row := []interface{}{
			winner.DrawTime.Format("2006-01-02 15:04:05"),
			winner.PrizeName,
			winner.WinnerID,
			winner.WinnerName,
			winner.PrizeLevel,
		}
		cell := fmt.Sprintf("A%d", startRow+i)
		if err := f.SetSheetRow(SheetWinners, cell, &row); err != nil {
			return fmt.Errorf("failed to write winner row %d: %w", i, err)
		}
	}

	// Save file
	if err := f.Save(); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	fmt.Printf("Successfully saved %d winners to Excel file [%s]\n", len(winners), filePath)
	return nil
}

// CreateExcelTemplate creates a new Excel template with sample data
func CreateExcelTemplate(filePath string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close Excel file: %v\n", err)
		}
	}()

	// Create Prizes sheet
	if err := f.SetSheetName("Sheet1", SheetPrizes); err != nil {
		return fmt.Errorf("failed to rename sheet: %w", err)
	}

	// Prizes header
	prizesHeader := []interface{}{"ID", "Name (CN)", "Name (EN)", "Count", "Level", "Probability"}
	if err := f.SetSheetRow(SheetPrizes, "A1", &prizesHeader); err != nil {
		return fmt.Errorf("failed to write prizes header: %w", err)
	}

	// Sample prizes data
	samplePrizes := [][]interface{}{
		{1, "特等奖：欧洲豪华双人游", "Grand Prize: Europe Luxury Tour", 1, 0, 0.1},
		{2, "一等奖：最新款笔记本电脑", "First Prize: Latest Laptop", 3, 1, 0.3},
		{3, "二等奖：降噪耳机", "Second Prize: Noise-Canceling Headphones", 5, 2, 0.6},
		{4, "三等奖：阳光普照购物卡", "Third Prize: Shopping Card", 10, 3, 0.9},
	}
	for i, prize := range samplePrizes {
		cell := fmt.Sprintf("A%d", i+2)
		if err := f.SetSheetRow(SheetPrizes, cell, &prize); err != nil {
			return fmt.Errorf("failed to write prize row: %w", err)
		}
	}

	// Create Participants sheet
	if _, err := f.NewSheet(SheetParticipants); err != nil {
		return fmt.Errorf("failed to create Participants sheet: %w", err)
	}

	// Participants header
	participantsHeader := []interface{}{"ID", "Name", "Department", "Email"}
	if err := f.SetSheetRow(SheetParticipants, "A1", &participantsHeader); err != nil {
		return fmt.Errorf("failed to write participants header: %w", err)
	}

	// Sample participants data
	sampleParticipants := [][]interface{}{
		{1, "张三", "技术部", "zhangsan@company.com"},
		{2, "李四", "市场部", "lisi@company.com"},
		{3, "王五", "人力资源部", "wangwu@company.com"},
		{4, "赵六", "财务部", "zhaoliu@company.com"},
		{5, "孙七", "技术部", "sunqi@company.com"},
		{6, "周八", "市场部", "zhouba@company.com"},
		{7, "吴九", "运营部", "wujiu@company.com"},
		{8, "郑十", "技术部", "zhengshi@company.com"},
		{9, "冯十一", "人力资源部", "fengshiyi@company.com"},
		{10, "陈十二", "财务部", "chenshier@company.com"},
	}
	for i, participant := range sampleParticipants {
		cell := fmt.Sprintf("A%d", i+2)
		if err := f.SetSheetRow(SheetParticipants, cell, &participant); err != nil {
			return fmt.Errorf("failed to write participant row: %w", err)
		}
	}

	// Create Winners sheet
	if _, err := f.NewSheet(SheetWinners); err != nil {
		return fmt.Errorf("failed to create Winners sheet: %w", err)
	}

	// Winners header
	winnersHeader := []interface{}{"Draw Time", "Prize Name", "Winner ID", "Winner Name", "Prize Level"}
	if err := f.SetSheetRow(SheetWinners, "A1", &winnersHeader); err != nil {
		return fmt.Errorf("failed to write winners header: %w", err)
	}

	// Set Prizes as the active sheet
	f.SetActiveSheet(0)

	// Save file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel template: %w", err)
	}

	fmt.Printf("Successfully created Excel template: %s\n", filePath)
	return nil
}
