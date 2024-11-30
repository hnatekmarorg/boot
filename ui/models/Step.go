package models

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

const (
	tickMark  = "✓"
	crossMark = "✗"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

// Step Allows to create multiple boolean tasks.
// If false it shows red X - task
// If it runs it shows spinner - task
// If it succeeds it shows green tick - task
type Step struct {
	Status  chan bool
	Purpose string // Purpose of step e.g. Installing something
	result  bool
	Spinner spinner.Model
	done    bool // If true it means that the task is done and we should not update it anymore.
}

func InitStep(purpose string, readChan chan bool) Step {
	spinnerComponent := spinner.New()
	spinnerComponent.Spinner = spinner.Dot
	spinnerComponent.Style = pendingStyle
	return Step{
		Spinner: spinnerComponent,
		Purpose: purpose,
		Status:  readChan,
		result:  false,
		done:    false,
	}
}

func (s Step) View() string {
	prefix := ""
	style := pendingStyle
	if !s.done {
		prefix = s.Spinner.View()
	} else {
		if s.done && s.result {
			style = successStyle
			prefix = style.Render(tickMark)
		} else {
			style = errorStyle
			prefix = style.Render(crossMark)
		}
	}
	return prefix + " " + s.Purpose
}

func (s Step) Init() tea.Cmd {
	if s.done {
		return nil
	}
	return s.Spinner.Tick
}

func (s Step) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if !s.done {
		select {
		case result := <-s.Status:
			s.result = result
			s.done = true
		case <-time.After(20 * time.Millisecond):
			s.Spinner, cmd = s.Spinner.Update(msg)
		}
	}
	return s, cmd
}
