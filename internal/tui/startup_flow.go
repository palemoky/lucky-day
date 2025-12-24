package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/palemoky/lucky-day/internal/i18n"
)

// StartupFlow represents the entire startup flow
type StartupFlow struct {
	stage      int // 0=language, 1=mode
	langModel  LanguageSelectionModel
	modeModel  ModeSelectionModel
	translator *i18n.Translator

	// Results
	selectedLang i18n.Language
	selectedMode LotteryMode
	userQuit     bool
}

// NewStartupFlow creates a new startup flow
func NewStartupFlow() StartupFlow {
	return StartupFlow{
		stage:     0,
		langModel: NewLanguageSelectionModel(),
	}
}

func (m StartupFlow) Init() tea.Cmd {
	return m.langModel.Init()
}

func (m StartupFlow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size for all stages
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		switch m.stage {
		case 0:
			langModel, _ := m.langModel.Update(wsMsg)
			m.langModel = langModel.(LanguageSelectionModel)
		case 1:
			modeModel, _ := m.modeModel.Update(wsMsg)
			m.modeModel = modeModel.(ModeSelectionModel)
		}
		return m, nil
	}

	switch m.stage {
	case 0: // Language selection
		var cmd tea.Cmd
		langModel, cmd := m.langModel.Update(msg)
		m.langModel = langModel.(LanguageSelectionModel)

		// Check if language selection is done
		if m.langModel.done {
			m.selectedLang = m.langModel.GetSelectedLanguage()
			m.translator = i18n.NewTranslator(m.selectedLang)
			m.stage = 1
			m.modeModel = NewModeSelectionModel(m.translator)
			// Copy window size from language model to mode model
			m.modeModel.width = m.langModel.width
			m.modeModel.height = m.langModel.height
			return m, m.modeModel.Init()
		}

		// Check if user quit
		if msg, ok := msg.(tea.KeyMsg); ok && (msg.String() == "q" || msg.String() == "ctrl+c") {
			if !m.langModel.done {
				m.userQuit = true
				return m, tea.Quit
			}
		}

		return m, cmd

	case 1: // Mode selection
		var cmd tea.Cmd
		modeModel, cmd := m.modeModel.Update(msg)
		m.modeModel = modeModel.(ModeSelectionModel)

		// Check if mode selection is done
		if m.modeModel.done {
			m.selectedMode = m.modeModel.GetSelectedMode()

			// Always quit - main will handle QR mode separately
			return m, tea.Quit
		}

		// Check if user quit
		if msg, ok := msg.(tea.KeyMsg); ok && (msg.String() == "q" || msg.String() == "ctrl+c") {
			if !m.modeModel.done {
				m.userQuit = true
				return m, tea.Quit
			}
		}

		return m, cmd
	}

	return m, nil
}

func (m StartupFlow) View() string {
	switch m.stage {
	case 0:
		return m.langModel.View()
	case 1:
		return m.modeModel.View()
	}
	return ""
}

// GetResults returns the selected language, mode, and quit status
func (m StartupFlow) GetResults() (i18n.Language, LotteryMode, bool) {
	return m.selectedLang, m.selectedMode, m.userQuit
}

// RunStartupFlow runs the unified startup flow
func RunStartupFlow() (i18n.Language, LotteryMode, bool, error) {
	m := NewStartupFlow()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return i18n.Chinese, ModeExcel, false, err
	}

	if flow, ok := finalModel.(StartupFlow); ok {
		lang, mode, quit := flow.GetResults()
		return lang, mode, quit, nil
	}

	return i18n.Chinese, ModeExcel, true, nil
}
