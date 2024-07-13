package tui

import (
	"essh/cmd/internal/session"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type editModel struct {
	session     session.Session
	updateState func(state) tea.Cmd
}

func initEditModel(session session.Session, updateState func(state) tea.Cmd) tea.Model {
	return editModel{
		session:     session,
		updateState: updateState,
	}
}

func (e editModel) Init() tea.Cmd {
	return nil
}

func (e editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			newState := state{
				page: listPage,
			}
			return e, e.updateState(newState)
		}
	}
	return e, nil
}

func (e editModel) View() string {
	return fmt.Sprintf("SessionName: %s\nUserName: %s\nHost: %s\nPort: %d\nPassword: %s\n",
		e.session.SessionName, e.session.UserName, e.session.Host, e.session.Port, e.session.Password)
}
