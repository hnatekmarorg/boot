package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hnatekmarorg/boot/ui/models"
	"math/rand/v2"
	"os"
	"time"
)

func main() {
	ch := make(chan bool, 1)
	step := models.Init("Sun", ch)
	go func() {
		time.Sleep(time.Second * 12)
		ch <- rand.IntN(1) == 0
	}()
	program := tea.NewProgram(step)
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
