package tui

import "essh/cmd/internal/session"

type state struct {
	page            int
	selectedIdx     int
	selectedSession session.Session
}
