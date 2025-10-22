package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ctrixcode/Pokemon-fetch-cli-Go/pokemon"
	"github.com/ctrixcode/Pokemon-fetch-cli-Go/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var fetchFlag = flag.Bool("fetch", false, "Fetch Pokemon data from API")
	var listFlag = flag.Bool("list", false, "Show interactive Pokemon list")
	flag.Parse()

	// If no flags are provided, default to fetch behavior for backward compatibility
	if !*fetchFlag && !*listFlag {
		fmt.Println("Fetching Pokemon data...")
		pokemon.StoreData()
		fmt.Println("Pokemon data fetched successfully! Use --list to view the interactive list.")
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
