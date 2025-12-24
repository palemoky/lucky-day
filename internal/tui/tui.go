package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/palemoky/lucky-day/internal/lottery"
	model1 "github.com/palemoky/lucky-day/internal/model"
)

// 最大显示抽奖人数
const maxDisplayedWinners = 5

// 定义TUI的几种状态
type appState int

const (
	statePrizeSelection appState = iota // 奖项选择
	stateDrawing                        // 正在抽奖（动画）
	stateShowWinners                    // 显示本次中奖结果
)

type model struct {
	engine         *lottery.Engine
	state          appState
	cursor         int // 当前选中的奖项索引
	width, height  int
	spinner        spinner.Model
	rollingNames   []string // 抽奖动画中滚动的名字
	currentWinners []model1.Participant
	lastErr        string
}

// NewTUIModel 创建并初始化一个新的TUI模型
func NewTUIModel(engine *lottery.Engine) *model {
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &model{
		engine:  engine,
		state:   statePrizeSelection,
		spinner: s,
	}
}

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	// 只在抽奖状态下更新动画
	if m.state == stateDrawing {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		prize := m.engine.GetPrizes()[m.cursor]
		count := prize.Count - prize.DrawnCount
		m.rollingNames = m.engine.GetRandomNames(count)

		return m, tea.Batch(cmd, tick())
	}

	return m, nil
}

func (m *model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case statePrizeSelection:
		return m.updatePrizeSelection(msg)
	case stateDrawing:
		return m.updateDrawing(msg)
	case stateShowWinners:
		return m.updateShowWinners(msg)
	}
	return m, nil
}

// 处理奖项选择界面的按键
func (m *model) updatePrizeSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prizes := m.engine.GetPrizes()
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(prizes)-1 {
			m.cursor++
		}
	case "enter", " ":
		prize := prizes[m.cursor]
		if prize.DrawnCount >= prize.Count {
			m.lastErr = fmt.Sprintf("错误: [%s] 的名额已抽完！", prize.Name)
			return m, nil
		}
		m.lastErr = ""
		m.state = stateDrawing
		return m, tick()
	case "r":
		prizeToReset := prizes[m.cursor]
		m.engine.ResetPrize(prizeToReset.ID)
		m.lastErr = fmt.Sprintf("提示: [%s] 已重置。", prizeToReset.Name)
	}
	return m, nil
}

// 处理抽奖动画界面的按键
func (m *model) updateDrawing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	default:
		prize := m.engine.GetPrizes()[m.cursor]
		winners, ok := m.engine.Draw(prize.ID)
		if !ok {
			m.lastErr = "抽奖失败，可能没有足够的候选人。"
		}
		m.currentWinners = winners
		m.state = stateShowWinners
		return m, nil
	}
}

// 处理显示中奖者界面的按键
func (m *model) updateShowWinners(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	default:
		m.currentWinners = nil
		m.state = statePrizeSelection
		return m, nil
	}
}

func (m *model) View() string {
	var mainContent, sidebar string

	// 根据状态渲染主面板
	switch m.state {
	case statePrizeSelection:
		mainContent = m.viewPrizeSelection()
	case stateDrawing:
		mainContent = m.viewDrawing()
	case stateShowWinners:
		mainContent = m.viewShowWinners()
	}

	// 渲染侧边栏
	sidebar = m.viewAllWinners()

	// 将主面板和侧边栏水平组合
	ui := lipgloss.JoinHorizontal(lipgloss.Top, mainContent, sidebar)

	// 构建最终视图：标题、UI、页脚垂直组合
	view := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("✨ Lucky Day ✨"),
		ui,
		m.viewFooter(),
	)

	// 将整个视图在屏幕上居中
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

// 渲染所有中奖者名单 (侧边栏)
func (m *model) viewAllWinners() string {
	var b strings.Builder
	b.WriteString("🏆 总中奖名单 🏆\n\n")

	allWinners := m.engine.GetAllWinners()
	prizes := m.engine.GetPrizes()

	hasPrintedFirstBlock := false
	for _, prize := range prizes {
		if winners, ok := allWinners[prize.ID]; ok && len(winners) > 0 {
			if hasPrintedFirstBlock {
				b.WriteString("\n")
			}

			b.WriteString(focusedStyle.Render(fmt.Sprintf("%s (%d/%d):", prize.Name, len(winners), prize.Count)))
			b.WriteString("\n")

			var names []string
			for _, w := range winners {
				names = append(names, w.Name)
			}

			// 每行最多显示5个名字
			const namesPerLine = 5
			for i := 0; i < len(names); i += namesPerLine {
				end := min(i+namesPerLine, len(names))
				line := strings.Join(names[i:end], ", ")
				b.WriteString(winnerStyle.Render("  " + line))
				b.WriteString("\n")
			}
			hasPrintedFirstBlock = true
		}
	}

	if !hasPrintedFirstBlock {
		b.WriteString(blurredStyle.Render("  还未有人中奖..."))
	}

	return sidebarStyle.Render(strings.TrimSuffix(b.String(), "\n"))
}

// 渲染奖项选择列表
func (m *model) viewPrizeSelection() string {
	var s strings.Builder
	s.WriteString("请选择要抽取的奖项：\n\n")

	prizes := m.engine.GetPrizes()
	for i, p := range prizes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		status := fmt.Sprintf("(%d/%d)", p.DrawnCount, p.Count)
		line := fmt.Sprintf("%s %s %s", cursor, p.Name, status)
		if m.cursor == i {
			s.WriteString(focusedStyle.Render(line))
		} else {
			s.WriteString(blurredStyle.Render(line))
		}
		s.WriteString("\n")
	}

	if m.lastErr != "" {
		s.WriteString("\n" + errorStyle.Render(m.lastErr))
	}

	return mainPanelStyle.Render(strings.TrimSuffix(s.String(), "\n"))
}

// 渲染抽奖动画
func (m *model) viewDrawing() string {
	var s strings.Builder
	prize := m.engine.GetPrizes()[m.cursor]

	s.WriteString(fmt.Sprintf("正在抽取 [%s] ... %s\n\n", prize.Name, m.spinner.View()))

	var winnerBlocks []string
	for i, name := range m.rollingNames {
		if i >= maxDisplayedWinners {
			// 如果达到显示上限，则添加一个省略号块并停止
			ellipsis := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(1, 2).
				Render("...")
			winnerBlocks = append(winnerBlocks, ellipsis)
			break
		}
		winnerBlocks = append(winnerBlocks, winnerBoxStyle.Render(name))
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, winnerBlocks...))
	s.WriteString("\n\n按下 [任意键] 停止！")

	return mainPanelStyle.Render(s.String())
}

// 渲染本次中奖结果
func (m *model) viewShowWinners() string {
	var s strings.Builder
	prize := m.engine.GetPrizes()[m.cursor]

	if len(m.currentWinners) == 0 {
		s.WriteString(fmt.Sprintf("很遗憾，[%s] 本次无人中奖。\n", prize.Name))
	} else {
		s.WriteString(fmt.Sprintf("🎉 恭喜以下人员获得 [%s] 🎉\n\n", prize.Name))
		var winnerBlocks []string
		for i, w := range m.currentWinners {
			if i >= maxDisplayedWinners {
				ellipsis := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(1, 2).Render("...")
				winnerBlocks = append(winnerBlocks, ellipsis)
				break
			}
			winnerBlocks = append(winnerBlocks, winnerBoxStyle.Render(w.Name))
		}
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, winnerBlocks...))
	}

	s.WriteString("\n\n按下 [任意键] 返回奖项列表。")

	return mainPanelStyle.Render(s.String())
}

// 渲染页脚
func (m *model) viewFooter() string {
	var instructions string
	switch m.state {
	case statePrizeSelection:
		instructions = "↑/↓: 选择 | Enter: 抽奖 | r: 重置当前奖项 | q: 退出"
	case stateDrawing:
		instructions = "任意键: 停止抽奖 | q: 退出"
	case stateShowWinners:
		instructions = "任意键: 返回 | q: 退出"
	}
	return helpStyle.Render("\n" + instructions)
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// StartTUI 启动TUI程序
func StartTUI(engine *lottery.Engine) error {
	p := tea.NewProgram(NewTUIModel(engine), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
