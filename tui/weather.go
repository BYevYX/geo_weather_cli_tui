package tui

import (
	"fmt"
	"geo-weather-cli/api"
	"reflect"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type weatherCurrentModel struct {
	locationInput textinput.Model
	done          bool
	weather       api.Current
}

func GetWeatherCurrentModel() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "London or 51.5074,0.1278"
	ti.Focus()
	ti.Width = 50

	return &weatherCurrentModel{
		locationInput: ti,
	}
}

func (m *weatherCurrentModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *weatherCurrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			coords := api.GetCoordinatesFromAddres(m.locationInput.Value()).Results[0]
			m.weather = api.GetCurrentWeather(coords.Lat, coords.Lon, coords.Timezone.Name).Current
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.locationInput, cmd = m.locationInput.Update(msg)
	return m, cmd
}

func (m *weatherCurrentModel) View() string {
	if m.done {
		val := reflect.ValueOf(m.weather)
		typ := reflect.TypeOf(m.weather)

		return getFinalString(val, typ)
	}
	return fmt.Sprintf(
		"Enter location (city or coordinates):\n\n%s\n\n(Enter to submit, Esc to quit)",
		m.locationInput.View(),
	)
}



type weatherForecastModel struct {
	inputs  []textinput.Model
	focused int
	weather api.Daily
	done    bool
}

func GetWeatherForecastModel() tea.Model {
	location := textinput.New()
	location.Placeholder = "New York"
	location.Focus()
	location.Width = 30

	days := textinput.New()
	days.Placeholder = "3"
	days.Width = 5

	return &weatherForecastModel{
		inputs: []textinput.Model{
			location,
			days,
		},
	}
}

func (m *weatherForecastModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *weatherForecastModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				coords := api.GetCoordinatesFromAddres(m.inputs[0].Value()).Results[0]
				m.weather = api.GetForecast(coords.Lat, coords.Lon, m.inputs[1].Value(), coords.Timezone.Name).Daily
				m.done = true
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *weatherForecastModel) View() string {
	if m.done {
		days, _ := strconv.Atoi(m.inputs[1].Value())

		var finalBuilder strings.Builder

		finalBuilder.WriteString(fmt.Sprintf(
			"\nGetting %d-day forecast for: %s\n",
			days,
			m.inputs[0].Value(),
		))

		val := reflect.ValueOf(m.weather)
		typ := reflect.TypeOf(m.weather)

		finalBuilder.WriteString(getFinalString(val, typ))

		return finalBuilder.String()
	}

	var b strings.Builder
	b.WriteString("Weather Forecast Request\n\n")
	b.WriteString("Location:\n")
	b.WriteString(m.inputs[0].View())
	b.WriteString("\n\nDays (1-10):\n")
	b.WriteString(m.inputs[1].View())
	b.WriteString("\n\n(Tab/Shift+Tab to switch fields, Enter to submit)")

	return b.String()
}

func (m *weatherForecastModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *weatherForecastModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
	m.updateFocus()
}

func (m *weatherForecastModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	m.updateFocus()
}

func (m *weatherForecastModel) updateFocus() {
    for i := 0; i < len(m.inputs); i++ {
        if i == m.focused {
            m.inputs[i].Focus()
        } else {
            m.inputs[i].Blur()
        }
    }
}

func getFinalString(val reflect.Value, typ reflect.Type) string {
	var finalBuilder strings.Builder

	for i := range val.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		
		if fieldValue.Kind() == reflect.Slice {
			var elements []string
			for j := 0; j < fieldValue.Len(); j++ {
				elements = append(elements, fmt.Sprintf("%v", fieldValue.Index(j).Interface()))
			}
			finalBuilder.WriteString(InputStyle.Render(
				fmt.Sprintf("%s: [%s]", jsonTag, strings.Join(elements, ", ")),
			))
			finalBuilder.WriteString("\n")
		} else {
			finalBuilder.WriteString(InputStyle.Render(
				fmt.Sprintf("%s: %v", jsonTag, fieldValue.Interface()),
			))
			finalBuilder.WriteString("\n")
		}
	}

	return finalBuilder.String()
}
