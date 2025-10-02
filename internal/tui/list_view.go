package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PokemonListItem struct {
	ID   int
	Name string
}

func (m Model) renderListView() string {
	var sb strings.Builder

	sb.WriteString("Pokemon List\n")
	sb.WriteString("============\n\n")

	pokemonList := m.getPokemonList()

	if len(pokemonList) == 0 {
		sb.WriteString("No Pokemon data found. Please ensure data is fetched first.\n")
	} else {
		sb.WriteString(fmt.Sprintf("Total Pokemon: %d\n\n", len(pokemonList)))

		displayCount := 20
		if len(pokemonList) < displayCount {
			displayCount = len(pokemonList)
		}

		for i := 0; i < displayCount; i++ {
			pokemon := pokemonList[i]
			sb.WriteString(fmt.Sprintf("%d. %s (ID: %d)\n", i+1, pokemon.Name, pokemon.ID))
		}

		if len(pokemonList) > displayCount {
			sb.WriteString(fmt.Sprintf("\n... and %d more\n", len(pokemonList)-displayCount))
		}
	}

	sb.WriteString("\nPress 'q' to quit\n")

	return sb.String()
}

func (m Model) getPokemonList() []PokemonListItem {
	var pokemonList []PokemonListItem

	dataDir := "data"

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return pokemonList
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		return pokemonList
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			idStr := strings.TrimSuffix(file.Name(), ".json")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			data, err := os.ReadFile(filepath.Join(dataDir, file.Name()))
			if err != nil {
				continue
			}

			var pokemonData struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}

			if err := json.Unmarshal(data, &pokemonData); err != nil {
				continue
			}

			pokemonList = append(pokemonList, PokemonListItem{
				ID:   id,
				Name: pokemonData.Name,
			})
		}
	}

	sort.Slice(pokemonList, func(i, j int) bool {
		return pokemonList[i].ID < pokemonList[j].ID
	})

	return pokemonList
}
