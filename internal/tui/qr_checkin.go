package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/palemoky/lucky-day/internal/checkin"
	"github.com/palemoky/lucky-day/internal/i18n"
)

// QRCheckInModel represents the QR check-in waiting screen
type QRCheckInModel struct {
	server     *checkin.Server
	translator *i18n.Translator
	qrPath     string
	url        string
	done       bool
	width      int
	height     int
}

// NewQRCheckInModel creates a new QR check-in model
func NewQRCheckInModel(server *checkin.Server, translator *i18n.Translator, qrPath, url string) QRCheckInModel {
	return QRCheckInModel{
		server:     server,
		translator: translator,
		qrPath:     qrPath,
		url:        url,
		done:       false,
	}
}

type participantCountMsg int

func (m QRCheckInModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		checkParticipantCount(m.server),
	)
}

func (m QRCheckInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.done = false
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "s":
			// Save participants to Excel
			// TODO: Implement save functionality
			return m, nil
		}

	case participantCountMsg:
		// Update participant count and schedule next check
		return m, checkParticipantCount(m.server)
	}

	return m, nil
}

func (m QRCheckInModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Align(lipgloss.Center).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Align(lipgloss.Center)

	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(2)

	// Build the content
	var content string

	// Title
	content += titleStyle.Render(m.translator.T("qr.title")) + "\n\n"

	// Instructions
	content += boxStyle.Render(m.translator.T("qr.instruction")) + "\n\n"

	// URL
	content += fmt.Sprintf("%s:\n", m.translator.T("qr.url"))
	content += urlStyle.Render(m.url) + "\n\n"

	// QR Code file info
	content += instructionStyle.Render(fmt.Sprintf("QR Code: %s", m.qrPath)) + "\n\n"

	// Participant count
	count := m.server.GetParticipantCount()
	content += fmt.Sprintf("%s: ", m.translator.T("qr.count"))
	content += countStyle.Render(fmt.Sprintf("%d", count)) + "\n\n"

	// Help text
	helpText := fmt.Sprintf(
		"%s | %s | %s",
		m.translator.T("qr.start"),
		m.translator.T("qr.save"),
		m.translator.T("qr.quit"),
	)
	content += helpStyle.Render(helpText)

	// Center the content on screen
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// IsDone returns whether the user pressed Enter to start lottery
func (m QRCheckInModel) IsDone() bool {
	return m.done
}

// checkParticipantCount returns a command that checks participant count periodically
func checkParticipantCount(server *checkin.Server) tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return participantCountMsg(server.GetParticipantCount())
	})
}

// ShowQRCheckIn displays the QR check-in screen and waits for user to start lottery
func ShowQRCheckIn(server *checkin.Server, translator *i18n.Translator, qrPath, url string) (bool, error) {
	m := NewQRCheckInModel(server, translator, qrPath, url)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	if qrModel, ok := finalModel.(QRCheckInModel); ok {
		return qrModel.IsDone(), nil
	}

	return false, nil
}
