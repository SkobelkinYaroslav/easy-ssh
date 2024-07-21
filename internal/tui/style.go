package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f7f7f7"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	cursorStyle  = focusedStyle
)
