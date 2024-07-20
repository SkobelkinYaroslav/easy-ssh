package tui

import (
	"essh/internal/session"
	client "essh/internal/ssh"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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
			return l, updateListItemFunc(len(l.connections), session.New("", "", "", "", 0))
		case "c":
			if len(l.connections) == 0 {
				return l, nil
			}
			conn, err := client.ConnectWithPassword(l.connections[l.cursor])
			if err != nil {
				l.connections[l.cursor].SetConnectable(false)
				return l, nil
			}
			err = client.SpawnShell(conn)
			if err != nil {
				l.connections[l.cursor].SetConnectable(false)
				return l, nil
			}
			return l, nil

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
	if len(l.connections) == 0 {
		return "No connections available.\n"
	}

	maxLength := 0
	for _, conn := range l.connections {
		displayName := conn.SessionName
		if !conn.IsConnectable {
			displayName = errorStyle.Render(conn.SessionName + " (unreachable)")
		} else {
			displayName = focusedStyle.Render(conn.SessionName)
		}
		if len(displayName) > maxLength {
			maxLength = len(displayName)
		}
	}

	s := "Current connections:\n"
	for i, connection := range l.connections {
		prefix := "  "
		if l.cursor == i {
			prefix = "> "
		}

		sessionName := connection.SessionName
		if !connection.IsConnectable {
			sessionName += " (unreachable)"
			sessionName = errorStyle.Render(sessionName)
		} else {
			if l.cursor == i {
				sessionName = focusedStyle.Render(sessionName)
			} else {
				sessionName = blurredStyle.Render(sessionName)
			}
		}

		padding := strings.Repeat(" ", maxLength-len(sessionName))
		s += fmt.Sprintf("%s%s%s\n", prefix, sessionName, padding)
	}
	return s
}
