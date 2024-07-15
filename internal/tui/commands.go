package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func updateListFunc() tea.Cmd {
	return func() tea.Msg {
		return updateList{}
	}
}

func openEditorFunc(ind int) tea.Cmd {
	return func() tea.Msg {
		return openEditor{selectedSession: ind}
	}
}
