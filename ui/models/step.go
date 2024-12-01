package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

const (
	tickMark  = "✓"
	crossMark = "✗"
	hourglass = "⌛"
)

type TaskResponse struct {
	Status  bool
	Message string
}

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
	Status  chan TaskResponse
	Purpose string // Purpose of step e.g. Installing something
	result  TaskResponse
	Spinner spinner.Model
	done    bool // If true it means that the task is done and we should not update it anymore.
}

func InitStep(purpose string, readChan chan TaskResponse) Step {
	spinnerComponent := spinner.New()
	spinnerComponent.Spinner = spinner.Dot
	spinnerComponent.Style = pendingStyle
	return Step{
		Spinner: spinnerComponent,
		Purpose: purpose,
		Status:  readChan,
		result:  TaskResponse{},
		done:    false,
	}
}

func (s Step) View() string {
	prefix := ""
	style := pendingStyle
	if !s.done {
		prefix = hourglass
	} else {
		status := fmt.Sprintf(" (%s)", s.result.Message)
		if s.result.Status {
			style = successStyle
			prefix = style.Render(tickMark) + status
		} else {
			style = errorStyle
			prefix = style.Render(crossMark) + status
		}
	}
	return prefix + " " + s.Purpose
}

func (s Step) Init() tea.Cmd {
	return s.Spinner.Tick
}

func (s Step) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
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
