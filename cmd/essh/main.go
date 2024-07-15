package main

import (
	"essh/internal/session"
	"essh/internal/tui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	sessions := []session.Session{
		session.New("Session1", "User1", "host1.example.com", "password1", 22),
		session.New("Session2", "User2", "host2.example.com", "password2", 2222),
		session.New("Session3", "User3", "host3.example.com", "password3", 3333),
	}
	p := tea.NewProgram(tui.InitMainModel(sessions))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
