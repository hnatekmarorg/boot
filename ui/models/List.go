package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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
	var cmd tea.Cmd
	for i := range s.steps {
		model, stepCMD := s.steps[i].Update(msg)
		cmd = tea.Batch(cmd, stepCMD)
		s.steps[i] = model.(Step)
	}

	// Calculate the progress based on the number of completed steps
	completedSteps := 0
	for _, step := range s.steps {
		if step.done { // Assuming Step has a method IsCompleted to check if it's done
			completedSteps++
		}
	}
	totalSteps := len(s.steps)
	if totalSteps > 0 {
		s.progressBar.SetPercent(float64(completedSteps) / float64(totalSteps))
	}

	return s, cmd
}

func (s StepList) View() string {
	result := strings.Builder{}
	for i := range s.steps {
		result.WriteString(s.steps[i].View() + "\n")
	}
	result.WriteString(s.progressBar.View())
	return result.String()
}
