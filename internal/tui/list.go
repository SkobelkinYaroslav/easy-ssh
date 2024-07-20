package tui

import (
	"essh/internal/session"
	client "essh/internal/ssh"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

type listModel struct {
	cursor      int
	list        list.Model
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
	n := len(l.connections)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			l.cursor = (l.cursor - 1 + n) % n
		case "down":
			l.cursor = (l.cursor + 1) % n
		case "enter":
			return l, updateListItemFunc(l.cursor, l.connections[l.cursor])
		case "a":
			return l, updateListItemFunc(len(l.connections), session.NewDefault())
		case "c":
			conn, err := client.ConnectWithPassword(l.connections[l.cursor])
			if err != nil {
				log.Println(err)
			}
			err = client.SpawnShell(conn)
			if err != nil {
				log.Println(err)
			}
			return l, tea.Quit

		case "d", "del":
			if len(l.connections) > 0 {
				l.connections = append((l.connections)[:l.cursor], (l.connections)[l.cursor+1:]...)
				if l.cursor >= len(l.connections) && len(l.connections) > 0 {
					l.cursor--
				}
			}
			return l, updateListFunc(l.connections)
		}
	}
	return l, nil
}

func (l listModel) View() string {
	lipgloss.NewStyle().Margin(1, 2).Render(l.list.View())
}

//func (l listModel) View() string {
//	if len(l.connections) == 0 {
//		return "No connections available.\n"
//	}
//
//	s := "Current connections:\n"
//	for i, connection := range l.connections {
//		prefix := "  "
//		if l.cursor == i {
//			prefix = "> "
//		}
//		s += fmt.Sprintf("%s%s\n", prefix, connection.SessionName)
//	}
//	return s
//}
