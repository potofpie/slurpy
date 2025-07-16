package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#F25D94")
	secondaryColor = lipgloss.Color("#EEEEEE")
	accentColor    = lipgloss.Color("#FFE66D")
	errorColor     = lipgloss.Color("#FF6B6B")
	successColor   = lipgloss.Color("#4ECDC4")

	// Base styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Underline(true)

	subHeaderStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Panel styles
	listPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1).
			Width(50).
			Height(30)

	listPanelFocusedStyle = listPanelStyle.Copy().
				BorderForeground(accentColor)

	detailsPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(secondaryColor).
				Padding(1).
				Width(60).
				Height(30)

	detailsPanelFocusedStyle = detailsPanelStyle.Copy().
					BorderForeground(accentColor)

	// List styles
	paginationStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	// Status styles
	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	statusPendingStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)
)
