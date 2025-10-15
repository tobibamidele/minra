package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobibamidele/minra/internal/app"
)

func main()  {
	// Get the starting directory
	startDir := "."
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	// Create app
	application, err := app.New(startDir)
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		application,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running editor: %v\n", err)
		os.Exit(1)
	}
}