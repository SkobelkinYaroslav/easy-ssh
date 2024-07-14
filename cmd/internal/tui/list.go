package tui

import (
	"essh/cmd/internal/session"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	connections []session.Session
	updateState func(state) tea.Cmd
	curState    state
}

func initListModel(connections []session.Session, updateState func(state) tea.Cmd, curState state) tea.Model {
	return listModel{
		connections: connections,
		updateState: updateState,
		curState:    curState,
	}
}

func (l listModel) Init() tea.Cmd {
	return nil
}

func (l listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return l.handleKeyMsg(msg)
	}
	return l, nil
}

func (l listModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	n := len(l.connections)
	switch msg.String() {
	case "up":
		l.curState.SetCurIdx((l.curState.selectedIdx - 1 + n) % n)
	case "down":
		l.curState.SetCurIdx((l.curState.selectedIdx + 1) % n)
	case "enter":
		l.curState.SetPage(editPage)
		l.curState.SetSelectedSession(l.connections[l.curState.selectedIdx])
		return l, l.updateState(l.curState)
	}
	return l, nil
}

func (l listModel) View() string {
	s := "Current connections:\n"
	for i, connection := range l.connections {
		prefix := "  "
		if l.curState.selectedIdx == i {
			prefix = "> "
		}
		s += fmt.Sprintf("%s%s\n", prefix, connection.SessionName)
	}
	return s
}
