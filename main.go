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
	objectsInSpace := []string{
		"Earth", "Mars", "Venus", "Jupiter", "Saturn", "Uranus", "Neptune",
	}

	list := models.InitList()

	for i := range objectsInSpace {
		ch := make(chan bool, 1)
		list = list.AddStep(models.InitStep(objectsInSpace[i], ch))
		go func() {
			time.Sleep(time.Second * 1)
			ch <- rand.IntN(2) == 0
		}()
	}

	program := tea.NewProgram(list)
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
