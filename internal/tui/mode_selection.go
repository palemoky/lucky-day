package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/palemoky/lucky-day/internal/i18n"
)

// LotteryMode represents different lottery modes
type LotteryMode string

const (
	ModeExcel LotteryMode = "excel"
	ModeQR    LotteryMode = "qr"
	ModeDB    LotteryMode = "db"
)

// ModeSelectionModel represents the mode selection screen
type ModeSelectionModel struct {
	cursor     int
	choices    []LotteryMode
	selected   LotteryMode
	done       bool
	translator *i18n.Translator
	width      int
	height     int
}

// NewModeSelectionModel creates a new mode selection model
func NewModeSelectionModel(translator *i18n.Translator) ModeSelectionModel {
	return ModeSelectionModel{
		cursor:     0,
		choices:    []LotteryMode{ModeExcel, ModeQR, ModeDB},
		done:       false,
		translator: translator,
	}
}

func (m ModeSelectionModel) Init() tea.Cmd {
	return nil
}

func (m ModeSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ModeSelectionModel) View() tea.View {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginTop(2).
		MarginBottom(1)

	choiceStyle := lipgloss.NewStyle().
		PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		PaddingLeft(2)

	s := titleStyle.Render(m.translator.T("mode.select")) + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		name := ""
		switch choice {
		case ModeExcel:
			name = m.translator.T("mode.excel")
		case ModeQR:
			name = m.translator.T("mode.qr")
		case ModeDB:
			name = m.translator.T("mode.db")
		}

		if m.cursor == i {
			cursor = ">"
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n"
		} else {
			s += choiceStyle.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n"
		}
	}

	s += "\n" + lipgloss.NewStyle().Faint(true).Render(m.translator.T("mode.instruction"))

	// Use dynamic window size for centering, like the lottery interface
	v := tea.NewView(lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		s,
	))
	v.AltScreen = true
	return v
}

// GetSelectedMode returns the selected mode
func (m ModeSelectionModel) GetSelectedMode() LotteryMode {
	return m.selected
}

// SelectMode shows mode selection screen and returns the selected mode
func SelectMode(translator *i18n.Translator) (LotteryMode, bool, error) {
	m := NewModeSelectionModel(translator)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return ModeExcel, false, err
	}

	if modeModel, ok := finalModel.(ModeSelectionModel); ok {
		if modeModel.done {
			return modeModel.GetSelectedMode(), false, nil
		}
		// User pressed q to quit
		return ModeExcel, true, nil
	}

	// Default to Excel if something went wrong
	return ModeExcel, true, nil
}
