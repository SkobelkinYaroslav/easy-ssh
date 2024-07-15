package tui

import (
	"essh/cmd/internal/session"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type listModel struct {
	cursor      int
	connections []session.Session
}

func initListModel(connections []session.Session) tea.Model {
	return listModel{
		connections: connections,
	}
}

func (l listModel) Init() tea.Cmd {
	return nil
}

func (l listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return l.handleKeyMsg(msg)
	default:
		log.Println("list.go: unhandled message type ", msg)

	}
	return l, nil
}

func (l listModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	n := len(l.connections)
	switch msg.String() {
	case "up":
		l.cursor = (l.cursor - 1 + n) % n
		return l, nil
	case "down":
		l.cursor = (l.cursor + 1) % n
		return l, nil
	case "enter":
		return l, openEditorFunc(l.cursor)
	}
	return l, nil
}

func (l listModel) View() string {
	s := "Current connections:\n"
	for i, connection := range l.connections {
		prefix := "  "
		if l.cursor == i {
			prefix = "> "
		}
		s += fmt.Sprintf("%s%s\n", prefix, connection.SessionName)
	}
	return s
}
