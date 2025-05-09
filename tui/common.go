package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func RunBubbleTea(modelFunc func() tea.Model) cli.ActionFunc {
	return func(c *cli.Context) error {
		p := tea.NewProgram(modelFunc())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	}
}

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	InputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	ContinueStyle = lipgloss.NewStyle().Foreground(darkGray)
)


type tickMsg time.Time

type ProgressModel struct {
	progress   progress.Model
	onComplete tea.Cmd // Команда при завершении
	tickCmd    tea.Cmd // Команда анимации
	done       bool    // Флаг завершения
	width      int     // Ширина прогресса
}

// Init implements tea.Model.
func (m *ProgressModel) Init() tea.Cmd {
	return nil
}

func NewProgress(width int) *ProgressModel {
	pm := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
	)

	return &ProgressModel{
		progress: pm,
		tickCmd:  tickCmd(),
		width:    width,
	}
}

func (m *ProgressModel) OnComplete(cmd tea.Cmd) *ProgressModel {
	m.onComplete = cmd
	return m
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		return m, nil

	case tickMsg:
		if m.progress.Percent() >= 1.0 {
			m.done = true
			return m, m.onComplete
		}
		cmd := m.progress.IncrPercent(0.2)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *ProgressModel) View() string {
	if m.done {
		return ""
	}
	return "\n  " + m.progress.View() + "\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
