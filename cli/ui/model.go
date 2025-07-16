package ui

import (
	"fmt"
	"strings"

	"github.com/bobby/slurpy/pkg/models"
	"github.com/bobby/slurpy/pkg/storage"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines the key bindings
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Help    key.Binding
	Quit    key.Binding
	Refresh key.Binding
	Clear   key.Binding
	Tab     key.Binding
}

// ShortHelp returns key help
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Refresh}
}

// FullHelp returns full key help
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Tab, k.Refresh, k.Clear},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear namespace"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
}

// Model represents the application state
type Model struct {
	requests     []*models.LoggedRequest
	namespaces   []string
	currentNS    string
	list         list.Model
	storage      *storage.Storage
	width        int
	height       int
	focusedPanel int // 0 = list, 1 = details
	showHelp     bool
	err          error
}

// requestItem wraps LoggedRequest for the list component
type requestItem struct {
	*models.LoggedRequest
}

func (i requestItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.Method, i.URL, i.Namespace)
}

func (i requestItem) Title() string {
	status := "PENDING"
	if i.Response != nil {
		status = fmt.Sprintf("%d", i.Response.StatusCode)
	} else if i.Error != "" {
		status = "ERROR"
	}

	return fmt.Sprintf("%s %s [%s]", i.Method, i.URL, status)
}

func (i requestItem) Description() string {
	duration := i.Duration.Truncate(1000000) // Truncate to milliseconds
	timeStr := i.Timestamp.Format("15:04:05")

	desc := fmt.Sprintf("%s • %s • %s", timeStr, duration, i.Namespace)
	if i.Error != "" {
		desc += " • " + i.Error
	}
	return desc
}

// InitialModel creates the initial application model
func InitialModel() Model {
	storage, err := storage.New()
	if err != nil {
		return Model{err: err}
	}

	// Create list
	items := []list.Item{}
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "HTTP Requests"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	model := Model{
		storage:      storage,
		list:         l,
		focusedPanel: 0,
		currentNS:    "all",
	}

	return model
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadRequestsCmd(m.storage, m.currentNS),
		loadNamespacesCmd(m.storage),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := m.width / 2
		listHeight := m.height - 2

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp

		case key.Matches(msg, keys.Refresh):
			return m, tea.Batch(
				loadRequestsCmd(m.storage, m.currentNS),
				loadNamespacesCmd(m.storage),
			)

		case key.Matches(msg, keys.Tab):
			if m.focusedPanel == 0 {
				m.focusedPanel = 1
			} else {
				m.focusedPanel = 0
			}

		case key.Matches(msg, keys.Clear):
			if m.currentNS != "all" && m.currentNS != "" {
				return m, clearNamespaceCmd(m.storage, m.currentNS)
			}
		}

	case requestsLoadedMsg:
		m.requests = msg.requests
		items := make([]list.Item, len(m.requests))
		for i, req := range m.requests {
			items[i] = requestItem{req}
		}
		m.list.SetItems(items)

	case namespacesLoadedMsg:
		m.namespaces = msg.namespaces

	case namespacesClearedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			return m, loadRequestsCmd(m.storage, m.currentNS)
		}

	case errMsg:
		m.err = msg.err
	}

	// Update list if focused
	if m.focusedPanel == 0 {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
	}

	// Left panel (requests list)
	listStyle := listPanelStyle
	if m.focusedPanel == 0 {
		listStyle = listPanelFocusedStyle
	}

	leftPanel := listStyle.Render(m.list.View())

	// Right panel (request details)
	detailsStyle := detailsPanelStyle
	if m.focusedPanel == 1 {
		detailsStyle = detailsPanelFocusedStyle
	}

	var rightPanel string
	if len(m.requests) > 0 && m.list.Index() < len(m.requests) {
		rightPanel = detailsStyle.Render(m.renderRequestDetails(m.requests[m.list.Index()]))
	} else {
		rightPanel = detailsStyle.Render("No request selected")
	}

	// Main layout
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Help view
	if m.showHelp {
		help := helpStyle.Render(m.renderHelp())
		return lipgloss.JoinVertical(lipgloss.Left, mainView, help)
	}

	return mainView
}

// renderRequestDetails renders the details panel content
func (m Model) renderRequestDetails(req *models.LoggedRequest) string {
	if req == nil {
		return "No request selected"
	}

	var b strings.Builder

	// Request info
	b.WriteString(headerStyle.Render("REQUEST"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Method: %s\n", req.Method))
	b.WriteString(fmt.Sprintf("URL: %s\n", req.URL))
	b.WriteString(fmt.Sprintf("Timestamp: %s\n", req.Timestamp.Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("Duration: %v\n", req.Duration))
	b.WriteString(fmt.Sprintf("Namespace: %s\n", req.Namespace))

	if req.Error != "" {
		b.WriteString(fmt.Sprintf("Error: %s\n", req.Error))
	}

	// Request headers
	if len(req.Headers) > 0 {
		b.WriteString("\n")
		b.WriteString(subHeaderStyle.Render("Request Headers:"))
		b.WriteString("\n")
		for k, v := range req.Headers {
			b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	// Request body
	if req.Body != "" {
		b.WriteString("\n")
		b.WriteString(subHeaderStyle.Render("Request Body:"))
		b.WriteString("\n")
		body := req.Body
		if len(body) > 200 {
			body = body[:200] + "..."
		}
		b.WriteString(body)
		b.WriteString("\n")
	}

	// Response info
	if req.Response != nil {
		resp := req.Response
		b.WriteString("\n")
		b.WriteString(headerStyle.Render("RESPONSE"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("Status: %d\n", resp.StatusCode))
		b.WriteString(fmt.Sprintf("Size: %d bytes\n", resp.Size))

		// Response headers
		if len(resp.Headers) > 0 {
			b.WriteString("\n")
			b.WriteString(subHeaderStyle.Render("Response Headers:"))
			b.WriteString("\n")
			for k, v := range resp.Headers {
				b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
			}
		}

		// Response body
		if resp.Body != "" {
			b.WriteString("\n")
			b.WriteString(subHeaderStyle.Render("Response Body:"))
			b.WriteString("\n")
			body := resp.Body
			if len(body) > 300 {
				body = body[:300] + "..."
			}
			b.WriteString(body)
		}
	}

	return b.String()
}

// renderHelp renders the help text
func (m Model) renderHelp() string {
	return `
SLURPY - HTTP Request Logger & Debugger

Key Bindings:
  ↑/k, ↓/j     Navigate up/down in request list
  ←/h, →/l     Navigate left/right (not implemented)
  tab          Switch focus between panels
  r            Refresh requests
  c            Clear current namespace (when not viewing all)
  ?            Toggle this help
  q/esc        Quit

Focus:
  Yellow border = focused panel
  White border  = unfocused panel

The left panel shows all HTTP requests, the right panel shows detailed
information about the selected request, similar to browser dev tools.
`
}
