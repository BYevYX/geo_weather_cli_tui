package tui

import (
	"fmt"
	"geo-weather-cli/api"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type coordinatesModel struct {
	addressInput textinput.Model
	progress     *ProgressModel // Встроенный компонент
	state        state          // Состояние: input | loading | result
	coordinates  struct{
		lat float64
		lon float64
	}
}

type state int

const (
	stateInput state = iota
	stateLoading
	stateResult
)

func GetCoordinatesModel() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Moscow, Red Square"
	ti.PlaceholderStyle = ContinueStyle
	ti.Focus()

	return &coordinatesModel{
		addressInput: ti,
		progress:     NewProgress(40).OnComplete(nil),
		state:        stateInput,
	}
}

func (m *coordinatesModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *coordinatesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateInput:
		return m.updateInput(msg)
	case stateLoading:
		return m.updateLoading(msg)
	case stateResult:
		return m.updateResult(msg)
	}
	return m, nil
}

func (m *coordinatesModel) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			m.state = stateLoading
			return m, tea.Batch(
				m.getCoordinates(), // Запускаем "запрос"
				m.progress.tickCmd,  // Запускаем прогресс
			)
		}
	}

	var cmd tea.Cmd
	m.addressInput, cmd = m.addressInput.Update(msg)
	return m, cmd
}

func (m *coordinatesModel) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Делегируем обработку прогрессу
	newProgress, cmd := m.progress.Update(msg)
	m.progress = newProgress.(*ProgressModel)

	if m.progress.done {
		m.state = stateResult
	}
	return m, cmd
}

func (m *coordinatesModel) updateResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}


func (m *coordinatesModel) getCoordinates() tea.Cmd {
	return func() tea.Msg {
		coords := api.GetCoordinatesFromAddres(m.addressInput.Value())
		
		m.coordinates.lat = coords.Results[0].Lat
		m.coordinates.lon = coords.Results[0].Lon

		return nil
	}
}

func (m *coordinatesModel) View() string {
	var sb strings.Builder

	switch m.state {
	case stateInput:
		sb.WriteString("Enter address:\n\n")
		sb.WriteString(m.addressInput.View())
		sb.WriteString(ContinueStyle.Render("\n\n(Enter to submit)"))

	case stateLoading:
		sb.WriteString(fmt.Sprintf("Requesting coordinates for %s...\n", m.addressInput.Value()))
		sb.WriteString(m.progress.View())

	case stateResult:
		sb.WriteString(InputStyle.Render(fmt.Sprintf("\nCoordinates: %f, %f", m.coordinates.lat, m.coordinates.lon)))
		sb.WriteString(InputStyle.Render("\n\n"))
		sb.WriteString(ContinueStyle.Render("Press any key to continue"))
	}

	return lipgloss.NewStyle().Width(50).Render(sb.String())
}
