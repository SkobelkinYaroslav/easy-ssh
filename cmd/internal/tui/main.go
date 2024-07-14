package tui

import (
	"essh/cmd/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listPage = iota
	editPage
)

type mainModel struct {
	curState    state
	updateState func(state) tea.Cmd
	pageModels  []tea.Model
	sessions    []session.Session
}

func InitMainModel(sessions []session.Session) tea.Model {
	updateState := func(s state) tea.Cmd {
		return func() tea.Msg {
			return s
		}
	}

	initialState := state{
		page:        listPage,
		selectedIdx: 0,
	}

	m := mainModel{
		curState:    initialState,
		sessions:    sessions,
		updateState: updateState,
	}

	m.pageModels = []tea.Model{
		initListModel(sessions, m.updateState, m.curState),
		initEditModel(sessions[initialState.selectedIdx], m.curState, m.updateState),
	}

	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.QuitMsg:
		if err := session.SaveToFile(m.sessions); err != nil {
			return nil, nil
		}
	case state:
		m.curState = msg
		if m.curState.page == editPage {
			m.pageModels[editPage] = initEditModel(m.sessions[m.curState.selectedIdx], m.curState, m.updateState)
		}
		if m.curState.page == listPage {
			m.sessions[m.curState.selectedIdx] = m.curState.selectedSession
			m.pageModels[listPage] = initListModel(m.sessions, m.updateState, m.curState)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			currentPageModel := m.pageModels[m.curState.page]
			newPageModel, cmd := currentPageModel.Update(msg)
			m.pageModels[m.curState.page] = newPageModel
			return m, cmd
		}

	}

	return m, nil
}

func (m mainModel) View() string {
	return m.pageModels[m.curState.page].View()
}
