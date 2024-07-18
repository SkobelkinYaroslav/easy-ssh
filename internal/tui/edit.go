package tui

import (
	"essh/internal/session"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

type editModel struct {
	// Current session
	ind     int
	session session.Session

	// TUI stuff
	focusIndex int
	inputs     []textinput.Model
}

func initEditModel(ind int, session session.Session) tea.Model {
	m := editModel{
		ind:     ind,
		session: session,
		inputs:  make([]textinput.Model, 5),
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

			m.session.SetSessionName(m.inputs[0].Value())
			m.session.SetUserName(m.inputs[1].Value())
			m.session.SetHost(m.inputs[2].Value())
			m.session.SetPort(port)
			m.session.SetPassword(m.inputs[4].Value())

			return m, updateListItemFunc(m.ind, m.session)

		case "up", "down":
			s := msg.String()

			if s == "up" {
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

			return m, tea.Batch(cmds...)

		}

	}

	newInput, cmd := m.inputs[m.focusIndex].Update(msg)
	m.inputs[m.focusIndex] = newInput
	return m, cmd
}

func (m editModel) View() string {
	var s string

	for i := range m.inputs {
		s += m.inputs[i].View()
		if i < len(m.inputs)-1 {
			s += "\n"
		}
	}

	s += "\n\nPress enter to submit"

	return s
}
