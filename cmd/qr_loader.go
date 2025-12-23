package main

import (
	"fmt"

	"github.com/palemoky/lucky-day/internal/checkin"
	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/datasource"
	"github.com/palemoky/lucky-day/internal/i18n"
	"github.com/palemoky/lucky-day/internal/model"
)

// loadFromQRCheckInContinuous starts QR check-in server in background
func loadFromQRCheckInContinuous(translator *i18n.Translator, lang i18n.Language) ([]model.Prize, []model.Participant, error) {
	// Load prizes from Excel (we still need prizes configuration)
	dsCfg, err := config.LoadDataSourceConfig(".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	excelPath := dsCfg.Excel.Path
	if excelPath == "" {
		excelPath = "examples/lottery_template.xlsx"
	}

	prizes, err := datasource.LoadPrizesFromExcel(excelPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load prizes: %w", err)
	}

	// Start check-in server in background
	server := checkin.NewServer(8888, translator)
	if err := server.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start check-in server: %w", err)
	}

	// Generate QR code
	qrPath := "checkin_qr.png"
	url := server.GetURL()
	if err := checkin.GenerateQRCode(url, qrPath); err != nil {
		_ = server.Stop() // Ignore error on cleanup
		return nil, nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Show simple message - QR code is ready
	fmt.Printf("\n✅ %s\n", translator.T("qr.ready"))
	fmt.Printf("📱 %s: %s\n", translator.T("qr.url"), url)
	fmt.Printf("🖼️  %s: %s\n\n", translator.T("qr.qr_file"), qrPath)
	fmt.Printf("💡 %s\n", translator.T("qr.hint"))
	fmt.Printf("⏎  %s\n\n", translator.T("qr.press_enter"))

	// Wait for user to press Enter
	_, _ = fmt.Scanln() // Ignore input

	// Stop server - no more check-ins allowed
	_ = server.Stop() // Ignore error on shutdown

	// Get participants
	participants := server.GetParticipants()

	if len(participants) == 0 {
		return nil, nil, fmt.Errorf("%s", translator.T("qr.no_participants"))
	}

	fmt.Printf("✅ %s: %d\n\n", translator.T("qr.total_participants"), len(participants))

	return prizes, participants, nil
}
