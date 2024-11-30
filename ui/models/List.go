package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type StepList struct {
	steps       []Step
	progressBar progress.Model
}

func (s StepList) Init() tea.Cmd {
	return nil
}

func InitList() StepList {
	return StepList{
		steps:       []Step{},
		progressBar: progress.New(progress.WithDefaultGradient()),
	}
}

func (s StepList) AddStep(step Step) StepList {
	s.steps = append(s.steps, step)
	return s
}

func (s StepList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return nil
	})
}

func (s StepList) View() string {
	result := ""
	for i := range s.steps {
		result += s.steps[i].View() + "\n"
	}
	result += s.progressBar.View()
	return result
}
