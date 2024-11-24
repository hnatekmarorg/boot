package models

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

/**
 * Allows to create multiple boolean tasks.
 * If false it shows red X - task
 * If it runs it shows spinner - task
 * If it succeds it shows green tick - task
 */
type Step struct {
	Status  chan bool
	Purpose string // Purpose of step e.g Installing something
	result  *bool
	Spinner spinner.Model
}

func Init(purpose string, readChan chan bool) Step {
	spinnerComponent := spinner.New()
	spinnerComponent.Spinner = spinner.Dot
	spinnerComponent.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Step{
		Spinner: spinnerComponent,
		Purpose: purpose,
		Status:  readChan,
		result:  nil,
	}
}

func (s Step) View() string {
	prefix := ""
	if s.result == nil {
		prefix = s.Spinner.View()
	} else {
		if *s.result {
			prefix = "✓"
		} else {
			prefix = "X"
		}
	}
	return prefix + " " + s.Purpose
}

func (s Step) Init() tea.Cmd {
	return s.Spinner.Tick
}

func (s Step) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if s.result == nil {
		select {
		case result := <-s.Status:
			s.result = &result
		case <-time.After(time.Millisecond):
			s.Spinner, cmd = s.Spinner.Update(msg)
		}
	}
	return s, cmd
}
