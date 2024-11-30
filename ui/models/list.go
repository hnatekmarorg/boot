package models

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type StepList struct {
	steps       []Step
	progressBar progress.Model
}
type tickMsg time.Time

func (s StepList) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func InitList() StepList {
	return StepList{
		steps:       []Step{},
		progressBar: progress.New(progress.WithDefaultScaledGradient()),
	}
}

func (s StepList) AddStep(step Step) StepList {
	s.steps = append(s.steps, step)
	return s
}

func (s StepList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Calculate the progress based on the number of completed steps
	completedSteps := 0
	for _, step := range s.steps {
		if step.done { // Assuming Step has a method IsCompleted to check if it's done
			completedSteps++
		}
	}
	totalSteps := len(s.steps)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s, tea.Quit
	case tickMsg:
		var cmd tea.Cmd = tickCmd()
		for i := range s.steps {
			model, stepCMD := s.steps[i].Update(msg)
			cmd = tea.Batch(cmd, stepCMD)
			s.steps[i] = model.(Step)
		}
		if completedSteps > 0 {
			return s, tea.Batch(s.progressBar.SetPercent(float64(completedSteps)/float64(totalSteps)), cmd)
		}
		return s, cmd
	case progress.FrameMsg:
		model, updateCMD := s.progressBar.Update(msg)
		s.progressBar = model.(progress.Model)
		return s, updateCMD
	default:
		return s, nil
	}
}

func (s StepList) View() string {
	result := strings.Builder{}
	for i := range s.steps {
		result.WriteString(s.steps[i].View() + "\n")
	}
	result.WriteString(s.progressBar.View() + "\n")
	result.WriteString(helpStyle("Press any key to quit"))
	return result.String()
}
