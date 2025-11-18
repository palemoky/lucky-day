package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	mainPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60).
			MarginRight(2)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Width(45)

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	winnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("228")).Bold(true)
	helpStyle    = blurredStyle.Padding(0, 1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Italic(true)

	winnerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(lipgloss.Color("228")).
			Foreground(lipgloss.Color("228")).
			Bold(true).
		// Padding(1, 3).
		Margin(0, 1)
)
