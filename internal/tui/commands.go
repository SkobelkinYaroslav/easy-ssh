package tui

import (
	"essh/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

type updateList struct {
	sessions []session.Session
}

type updateListItem struct {
	ind     int
	session session.Session
}

func updateListFunc(sessions []session.Session) tea.Cmd {
	return func() tea.Msg {
		return updateList{sessions: sessions}
	}
}

func updateListItemFunc(ind int, s session.Session) tea.Cmd {
	return func() tea.Msg {
		return updateListItem{
			ind:     ind,
			session: s,
		}
	}
}
