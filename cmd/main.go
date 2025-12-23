package main

import (
	"fmt"
	"log"
	"os"

	"github.com/palemoky/lucky-day/internal/checkin"
	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/datasource"
	"github.com/palemoky/lucky-day/internal/i18n"
	"github.com/palemoky/lucky-day/internal/lottery"
	"github.com/palemoky/lucky-day/internal/model"
	"github.com/palemoky/lucky-day/internal/tui"
)

func main() {
	// Unified startup flow - no screen flicker!
	selectedLang, selectedMode, quit, err := tui.RunStartupFlow()
	if err != nil {
		log.Fatalf("Startup failed: %v", err)
	}
	if quit {
		fmt.Println("Goodbye!")
		return
	}

	translator := i18n.NewTranslator(selectedLang)

	var participants []model.Participant
	var prizes []model.Prize

	// Step 3: Load data based on selected mode
	switch selectedMode {
	case tui.ModeExcel:
		// Load from Excel
		prizes, participants, err = loadFromExcel(translator)
		if err != nil {
			log.Fatalf("%s: %v", translator.T("data.load_failed"), err)
		}

	case tui.ModeQR:
		// QR Check-in mode - run QR UI in continuation
		prizes, participants, err = loadFromQRCheckInContinuous(translator)
		if err != nil {
			log.Fatalf("%s: %v", translator.T("data.load_failed"), err)
		}

	case tui.ModeDB:
		// Load from database
		prizes, participants, err = loadFromDatabase(translator)
		if err != nil {
			log.Fatalf("%s: %v", translator.T("data.load_failed"), err)
		}

	default:
		log.Fatalf("Unknown mode: %s", selectedMode)
	}

	if len(participants) == 0 {
		log.Fatal(translator.T("data.empty_list"))
	}

	fmt.Printf("%s %d %s\n", translator.T("data.load_success"), len(participants), translator.T("data.participants"))
	fmt.Printf("%s %d %s\n", translator.T("data.load_success"), len(prizes), translator.T("data.prizes"))

	// Step 4: Initialize lottery engine
	engine := lottery.NewEngine(participants, prizes)

	// Step 5: Start TUI
	if err := tui.StartTUI(engine); err != nil {
		fmt.Printf("%s: %v\n", translator.T("app.error"), err)
		os.Exit(1)
	}

	fmt.Println(translator.T("app.exit"))
}

// loadFromQRCheckInContinuous starts QR check-in server in background
func loadFromQRCheckInContinuous(translator *i18n.Translator) ([]model.Prize, []model.Participant, error) {
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

// loadFromExcel loads prizes and participants from Excel file
func loadFromExcel(translator *i18n.Translator) ([]model.Prize, []model.Participant, error) {
	fmt.Println(translator.T("data.source_excel"))

	// Load configuration
	dsCfg, err := config.LoadDataSourceConfig(".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override type to excel if not set
	if dsCfg.Type != "excel" {
		dsCfg.Type = "excel"
		if dsCfg.Excel.Path == "" {
			dsCfg.Excel.Path = "lottery_template.xlsx"
		}
	}

	// Load prizes from Excel
	prizes, err := datasource.LoadPrizesFromExcel(dsCfg.Excel.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load prizes: %w", err)
	}

	// Load participants from Excel
	participants, err := datasource.LoadParticipantsFromExcel(dsCfg.Excel.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load participants: %w", err)
	}

	return prizes, participants, nil
}

// loadFromQRCheckIn starts QR check-in server and collects participants
// loadFromDatabase loads prizes and participants from database
func loadFromDatabase(translator *i18n.Translator) ([]model.Prize, []model.Participant, error) {
	fmt.Println(translator.T("data.source_db"))

	// Load configuration
	dsCfg, err := config.LoadDataSourceConfig(".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Load prizes from config (YAML)
	prizes, err := config.LoadPrizes(".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load prizes: %w", err)
	}

	// Load participants from database
	participants, err := datasource.LoadParticipants(dsCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load participants: %w", err)
	}

	return prizes, participants, nil
}
