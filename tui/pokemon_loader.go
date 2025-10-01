package tui

import (
	"ctrix/pokemon-fetch-cli-go/pokemon"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

// LoadPokemonFromFiles loads Pokemon data from JSON files and extracts names
func LoadPokemonFromFiles() tea.Msg {
	var pokemonList []PokemonListItem

	// Check if data directory exists
	dataDir := "data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return pokemonDataMsg{
			pokemonList: pokemonList,
			err:         fmt.Errorf("data directory not found - please run the fetch command first"),
		}
	}

	// Collect all JSON files first
	var files []string
	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return pokemonDataMsg{
			pokemonList: pokemonList,
			err:         fmt.Errorf("error reading data directory: %v", err),
		}
	}

	// Process each file and extract Pokemon data
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}

		var pokemonData pokemon.PokemonData
		if err := json.Unmarshal(data, &pokemonData); err != nil {
			continue // Skip files with invalid JSON
		}

		pokemonList = append(pokemonList, PokemonListItem{
			ID:   fmt.Sprintf("%d", pokemonData.Id),
			Name: strings.Title(pokemonData.Name),
			File: file,
		})
	}

	// Sort by Pokemon ID
	sort.Slice(pokemonList, func(i, j int) bool {
		id1, err1 := strconv.Atoi(pokemonList[i].ID)
		id2, err2 := strconv.Atoi(pokemonList[j].ID)
		if err1 != nil || err2 != nil {
			return pokemonList[i].ID < pokemonList[j].ID
		}
		return id1 < id2
	})

	return pokemonDataMsg{
		pokemonList: pokemonList,
		err:         nil,
	}
}