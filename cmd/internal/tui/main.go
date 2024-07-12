package tui

import (
	"essh/cmd/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	page        int
	selectedIdx int
}

const (
	listPage = iota
	editPage
)

type mainModel struct {
	curState    state
	pageModels  []tea.Model
	sessions    []session.Session
	updateState func(state) tea.Cmd
}

func InitMainModel() tea.Model {
	updateState := func(s state) tea.Cmd {
		return func() tea.Msg {
			return s
		}
	}
	sessions := []session.Session{
		session.New("Session1", "User1", "host1.example.com", "password1", 22),
		session.New("Session2", "User2", "host2.example.com", "password2", 2222),
		session.New("Session3", "User3", "host3.example.com", "password3", 3333),
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
		initListModel(sessions, m.updateState, &m.curState),
		initEditModel(sessions[initialState.selectedIdx], m.updateState),
	}

	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case state:
		m.curState = msg
		if m.curState.page == editPage {
			// Создаем новый editModel с выбранной сессией
			m.pageModels[editPage] = initEditModel(m.sessions[m.curState.selectedIdx], m.updateState)
		}
	case tea.KeyMsg:
		currentPageModel := m.pageModels[m.curState.page]
		newPageModel, cmd := currentPageModel.Update(msg)
		m.pageModels[m.curState.page] = newPageModel
		return m, cmd
	}

	return m, nil
}

func (m mainModel) View() string {
	return m.pageModels[m.curState.page].View()
}
