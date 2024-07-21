package tui

import (
	"essh/internal/session"
	client "essh/internal/ssh"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMapList struct {
	Up      key.Binding
	Down    key.Binding
	Edit    key.Binding
	Help    key.Binding
	Connect key.Binding
	Quit    key.Binding
	Delete  key.Binding
	Add     key.Binding
}

func (k keyMapList) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k keyMapList) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Help, k.Quit},
		{k.Add, k.Edit, k.Connect, k.Delete},
	}
}

type listModel struct {
	keys        keyMapList
	help        help.Model
	cursor      int
	connections []session.Session
}

func initListModel(connections []session.Session) tea.Model {
	var keys = keyMapList{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit entry"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Connect: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "connect"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete entry"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add entry"),
		),
	}

	return listModel{
		keys:        keys,
		help:        help.New(),
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
		case "?":
			l.help.ShowAll = !l.help.ShowAll
		}
	}
	return l, nil
}

func (l listModel) View() string {
	if len(l.connections) == 0 {
		return "No connections available.\n"
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

		s += fmt.Sprintf("%s%s\n", prefix, sessionName)
	}

	s += fmt.Sprintf("%s\n", l.help.View(l.keys))

	return s
}
