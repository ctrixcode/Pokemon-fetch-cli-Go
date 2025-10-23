package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type exportSuccessMsg struct {
	filePath string
}

type exportErrorMsg struct {
	err error
}

// ExportPokemonList exports the full Pokemon list to a JSON file
func ExportPokemonList(pokemons []PokemonListItem) tea.Cmd {
	return func() tea.Msg {
		timestamp := time.Now().Format("20060102_150405")
		filePath := fmt.Sprintf("pokemon_list_%s.json", timestamp)

		file, err := os.Create(filePath)
		if err != nil {
			return exportErrorMsg{err: err}
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(pokemons); err != nil {
			return exportErrorMsg{err: err}
		}

		return exportSuccessMsg{filePath: filePath}
	}
}

// ExportSinglePokemon exports a single Pokemon to a JSON file
func ExportSinglePokemon(pokemon PokemonListItem) tea.Cmd {
	return func() tea.Msg {
		timestamp := time.Now().Format("20060102_150405")
		filePath := fmt.Sprintf("pokemon_%s_%s.json", pokemon.ID, timestamp)

		file, err := os.Create(filePath)
		if err != nil {
			return exportErrorMsg{err: err}
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(pokemon); err != nil {
			return exportErrorMsg{err: err}
		}

		return exportSuccessMsg{filePath: filePath}
	}
}
