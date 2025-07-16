package ui

import (
	"github.com/bobby/slurpy/pkg/models"
	"github.com/bobby/slurpy/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Messages for async operations
type requestsLoadedMsg struct {
	requests []*models.LoggedRequest
}

type namespacesLoadedMsg struct {
	namespaces []string
}

type namespacesClearedMsg struct {
	err error
}

type errMsg struct {
	err error
}

// loadRequestsCmd loads requests for a namespace
func loadRequestsCmd(storage *storage.Storage, namespace string) tea.Cmd {
	return func() tea.Msg {
		var requests []*models.LoggedRequest
		var err error

		if namespace == "all" || namespace == "" {
			requests, err = storage.LoadAllRequests()
		} else {
			requests, err = storage.LoadRequests(namespace)
		}

		if err != nil {
			return errMsg{err}
		}

		return requestsLoadedMsg{requests}
	}
}

// loadNamespacesCmd loads all namespaces
func loadNamespacesCmd(storage *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		namespaces, err := storage.GetNamespaces()
		if err != nil {
			return errMsg{err}
		}

		return namespacesLoadedMsg{namespaces}
	}
}

// clearNamespaceCmd clears a namespace
func clearNamespaceCmd(storage *storage.Storage, namespace string) tea.Cmd {
	return func() tea.Msg {
		err := storage.ClearNamespace(namespace)
		return namespacesClearedMsg{err}
	}
}
