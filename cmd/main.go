package main

import (
	"ctrix/pokemon-fetch-cli-go/pokemon"
	"ctrix/pokemon-fetch-cli-go/tui"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var fetchFlag = flag.Bool("fetch", false, "Fetch Pokemon data from API")
	var listFlag = flag.Bool("list", false, "Show interactive Pokemon list")
	flag.Parse()

	// If no flags are provided, show help
	if !*fetchFlag && !*listFlag {
		showHelp()
		return
	}

	// Handle fetch command
	if *fetchFlag {
		fmt.Println("Fetching Pokemon data...")
		pokemon.StoreData()
		fmt.Println("Pokemon data fetched and stored successfully!")
		
		// If only fetch flag is provided, exit
		if !*listFlag {
			return
		}
	}

	// Handle list command
	if *listFlag {
		// Check if data exists
		if _, err := os.Stat("data"); os.IsNotExist(err) {
			fmt.Println("No Pokemon data found. Please run with --fetch first to download the data.")
			os.Exit(1)
		}

		// Start the TUI
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	}
}

func showHelp() {
	fmt.Println("Pokemon Fetch CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/main.go [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --fetch    Fetch Pokemon data from the API")
	fmt.Println("  --list     Show interactive Pokemon list (requires data to be fetched first)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/main.go --fetch           # Fetch data only")
	fmt.Println("  go run cmd/main.go --list            # Show list only (requires existing data)")
	fmt.Println("  go run cmd/main.go --fetch --list    # Fetch data and show list")
	fmt.Println()
}