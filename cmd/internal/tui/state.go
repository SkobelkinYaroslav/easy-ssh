package tui

import "essh/cmd/internal/session"

type state struct {
	page            int
	selectedIdx     int
	selectedSession session.Session
}

func NewState() *state {
	return &state{}
}

func (s *state) SetPage(n int) {
	s.page = n
}

func (s *state) SetCurIdx(n int) {
	s.selectedIdx = n
}

func (s *state) SetSelectedSession(selected session.Session) {
	s.selectedSession = selected
}
