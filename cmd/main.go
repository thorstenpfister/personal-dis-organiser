package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"personal-disorganizer/internal/app"
	"personal-disorganizer/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	purge := flag.Bool("purge", false, "Delete all data and start fresh")
	flag.Parse()

	// Handle purge command
	if *purge {
		if err := handlePurge(); err != nil {
			log.Fatalf("Failed to purge data: %v", err)
		}
		return
	}

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

func handlePurge() error {
	fmt.Print("Are you sure you want to delete all data? This cannot be undone. [Y/n]: ")
	
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" && response != "" {
		fmt.Println("Purge cancelled.")
		return nil
	}
	
	// Initialize storage to get the purge functionality
	storage, err := storage.NewStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	
	if err := storage.PurgeData(); err != nil {
		return fmt.Errorf("failed to purge data: %w", err)
	}
	
	fmt.Println("All data has been successfully deleted.")
	return nil
}