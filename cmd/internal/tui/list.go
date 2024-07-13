package tui

import (
	"essh/cmd/internal/session"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	cursor      int
	connections []session.Session
	updateState func(state) tea.Cmd
	curState    state
}

func initListModel(connections []session.Session, updateState func(state) tea.Cmd, curState state) tea.Model {
	return listModel{
		cursor:      0,
		connections: connections,
		updateState: updateState,
		curState:    curState,
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
			newState := l.curState
			newState.page = editPage
			newState.selectedIdx = l.cursor
			newState.selectedSession = l.connections[l.cursor]
			return l, l.updateState(newState)
		}
	}
	return l, nil
}

func (l listModel) View() string {
	s := "Current connections:\n"

	for i, connection := range l.connections {
		if l.cursor == i {
			s += "> "
		}
		s += fmt.Sprintf("%s\n", connection.SessionName)
	}

	return s
}
