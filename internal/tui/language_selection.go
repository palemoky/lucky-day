package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/palemoky/lucky-day/internal/i18n"
)

// LanguageSelectionModel represents the language selection screen
type LanguageSelectionModel struct {
	cursor   int
	choices  []i18n.Language
	selected i18n.Language
	done     bool
	width    int
	height   int
}

// NewLanguageSelectionModel creates a new language selection model
func NewLanguageSelectionModel() LanguageSelectionModel {
	return LanguageSelectionModel{
		cursor:  0,
		choices: []i18n.Language{i18n.Chinese, i18n.English},
		done:    false,
	}
}

func (m LanguageSelectionModel) Init() tea.Cmd {
	return nil
}

func (m LanguageSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
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

func (m LanguageSelectionModel) View() string {
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

	s := titleStyle.Render("请选择语言 / Select Language") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			name := ""
			if choice == i18n.Chinese {
				name = "中文"
			} else {
				name = "English"
			}
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n"
		} else {
			name := ""
			if choice == i18n.Chinese {
				name = "中文"
			} else {
				name = "English"
			}
			s += choiceStyle.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n"
		}
	}

	s += "\n" + lipgloss.NewStyle().Faint(true).Render("使用 ↑/↓ 选择，回车确认 | Use ↑/↓ to select, Enter to confirm")

	// Use dynamic window size for centering, like the lottery interface
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		s,
	)
}

// GetSelectedLanguage returns the selected language
func (m LanguageSelectionModel) GetSelectedLanguage() i18n.Language {
	return m.selected
}

// IsQuit returns whether the user quit without selecting
func (m LanguageSelectionModel) IsQuit() bool {
	return !m.done
}

// SelectLanguage shows language selection screen and returns the selected language
func SelectLanguage() (i18n.Language, bool, error) {
	m := NewLanguageSelectionModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return i18n.Chinese, false, err
	}

	if langModel, ok := finalModel.(LanguageSelectionModel); ok {
		if langModel.done {
			return langModel.GetSelectedLanguage(), false, nil
		}
		// User pressed q to quit
		return i18n.Chinese, true, nil
	}

	// Default to Chinese if something went wrong
	return i18n.Chinese, true, nil
}
