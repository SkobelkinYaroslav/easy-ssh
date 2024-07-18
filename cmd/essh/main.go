package main

import (
	"essh/internal/session"
	"essh/internal/tui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	sessions, err := session.GetFromFile()
	if err != nil {
		
	}
	p := tea.NewProgram(tui.InitMainModel(sessions))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
