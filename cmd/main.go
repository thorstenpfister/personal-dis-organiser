package main

import (
	"log"

	"personal-disorganizer/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the application model
	model, err := app.NewModel()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	// Create and run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}