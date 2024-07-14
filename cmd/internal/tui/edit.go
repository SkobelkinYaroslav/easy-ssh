package tui

import (
	"essh/cmd/internal/session"
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type editModel struct {
	// Manage state
	session     session.Session
	updateState func(state) tea.Cmd
	curState    state

	// TUI stuff
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initEditModel(session session.Session, curState state, updateState func(state) tea.Cmd) tea.Model {
	m := editModel{
		session:     session,
		updateState: updateState,
		curState:    curState,
		inputs:      make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Session Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.SetValue(session.SessionName)
		case 1:
			t.Placeholder = "User Name"
			t.SetValue(session.UserName)
		case 2:
			t.Placeholder = "Host"
			t.SetValue(session.Host)
		case 3:
			t.Placeholder = "Port"
			t.CharLimit = 6
			t.SetValue(strconv.Itoa(session.Port))
		case 4:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.SetValue("sample text")
		}

		m.inputs[i] = t
	}

	return m

}

func (m editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			port, err := strconv.Atoi(m.inputs[3].Value())
			if err != nil {
				return m, nil
			}
			newSession := session.Session{
				SessionName: m.inputs[0].Value(),
				UserName:    m.inputs[1].Value(),
				Host:        m.inputs[2].Value(),
				Port:        port,
				Password:    m.inputs[4].Value(),
			}

			m.curState.SetPage(listPage)
			m.curState.SetSelectedSession(newSession)
			m.curState.SetCurIdx(m.curState.selectedIdx)
			return m, m.updateState(m.curState)

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex = (m.focusIndex - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = blurredStyle
					m.inputs[i].TextStyle = blurredStyle
				}
			}
			cmds = append(cmds)

			return m, tea.Batch(cmds...)
		}
	}

	newInput, cmd := m.inputs[m.focusIndex].Update(msg)
	m.inputs[m.focusIndex] = newInput
	return m, cmd
}

func (m editModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\nPress enter to submit")

	return b.String()
}
