package main

import (
	"log"
	"os"

	"github.com/bobby/slurpy/cli/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
