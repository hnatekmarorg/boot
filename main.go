package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/hnatekmarorg/boot/ui/models"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
)

func findPort() (int, error) {
	// Scan ports from 8000 to end of range
	for port := 8000; port <= 65535; port++ {
		// Check if the port is available
		client, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			defer client.Close()
			return port, nil
		}
	}
	return 8080, fmt.Errorf("No available ports found")
}

var (
	detectPort      = make(chan models.TaskResponse)
	startServer     = make(chan models.TaskResponse)
	connectToRouter = make(chan models.TaskResponse)
)

func initSSHServer() {
	port, err := findPort()
	if err != nil {
		detectPort <- models.TaskResponse{
			Status:  true,
			Message: "Port not found!",
		}
		return
	}
	detectPort <- models.TaskResponse{
		Status:  true,
		Message: fmt.Sprintf("%d", port),
	}

	startServer <- models.TaskResponse{
		Status:  true,
		Message: "Server started!",
	}
	ssh.Handle(func(s ssh.Session) {
		cmd := exec.Command("/bin/bash")
		ptyReq, _, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				io.Copy(f, s) // stdin
			}()
			io.Copy(s, f) // stdout
			cmd.Wait()
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	err = ssh.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		startServer <- models.TaskResponse{
			Status:  false,
			Message: "Failed to start server",
		}
	}
}

func main() {
	steps := models.InitList()
	steps = steps.AddStep(models.InitStep("Find port", detectPort))
	steps = steps.AddStep(models.InitStep("Start server", startServer))
	steps = steps.AddStep(models.InitStep("Connect to router", connectToRouter))
	go initSSHServer()
	program := tea.NewProgram(steps)
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
