package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultTitle = "Unknown Request"
	defaultDesc  = "No details available"
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(requestItem)
	if !ok {
		return
	}

	title := i.Title()
	desc := i.Description()

	if title == "" {
		title = defaultTitle
	}
	if desc == "" {
		desc = defaultDesc
	}

	// Styling based on selection
	var titleStyle, descStyle lipgloss.Style

	if index == m.Index() {
		// Selected item
		titleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1)
		descStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Padding(0, 1)
	} else {
		// Normal item
		titleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Padding(0, 1)
		descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)
	}

	// Add status indicator
	status := "●"
	statusColor := lipgloss.Color("#888888")

	if i.Response != nil {
		if i.Response.StatusCode >= 200 && i.Response.StatusCode < 300 {
			statusColor = successColor
		} else if i.Response.StatusCode >= 400 {
			statusColor = errorColor
		} else {
			statusColor = accentColor
		}
	} else if i.Error != "" {
		statusColor = errorColor
	}

	statusStyle := lipgloss.NewStyle().Foreground(statusColor)

	// Truncate long titles and descriptions
	maxWidth := m.Width() - 6 // Account for padding and status indicator
	if len(title) > maxWidth {
		title = title[:maxWidth-3] + "..."
	}
	if len(desc) > maxWidth {
		desc = desc[:maxWidth-3] + "..."
	}

	titleLine := fmt.Sprintf("%s %s", statusStyle.Render(status), titleStyle.Render(title))
	descLine := descStyle.Render(desc)

	fmt.Fprint(w, titleLine)
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, descLine)
}
