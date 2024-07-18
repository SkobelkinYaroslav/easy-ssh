package tui

import (
	"essh/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listPage = iota
	editPage
)

type mainModel struct {
	page       int
	pageModels []tea.Model
	sessions   []session.Session
}

func InitMainModel(sessions []session.Session) tea.Model {
	m := mainModel{
		sessions: sessions,
	}

	m.pageModels = []tea.Model{
		initListModel(sessions),
		initEditModel(0, session.Session{}),
	}

	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if tea.Key(msg).String() == "ctrl+c" {
			err := session.SaveToFile(m.sessions)
			if err != nil {
				return nil, nil
			}
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.pageModels[m.page], cmd = m.pageModels[m.page].Update(msg)
		return m, cmd
	case updateList:
		l := initListModel(msg.sessions)
		m.sessions = msg.sessions
		m.page = listPage
		m.pageModels[listPage] = l
	case updateListItem:
		e := initEditModel(msg.ind, msg.session)
		m.page = editPage
		m.pageModels[editPage] = e
		//case updateListItem:
		//	m.sessions[msg.ind] = msg.session
		//	l := initListModel(m.sessions)
		//	m.page = listPage
		//	m.pageModels[listPage] = l
	}
	return m, nil
}
func (m mainModel) View() string {
	return m.pageModels[m.page].View()
}
