package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	listTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	listHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	listItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	listHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type PokemonListItem struct {
	ID   int
	Name string
}

func (m Model) renderListView() string {
	var sb strings.Builder

	// Clear any previous state and start fresh
	sb.WriteString(listTitleStyle.Render("📋 Pokémon List"))
	sb.WriteString("\n\n")

	pokemonList := m.getPokemonList()

	if len(pokemonList) == 0 {
		sb.WriteString("No Pokémon data found. Please ensure data is fetched first.\n\n")
		sb.WriteString(listHelpStyle.Render("Press 'l' for lookup • 'q' to quit"))
		return sb.String()
	}

	sb.WriteString(listHeaderStyle.Render(fmt.Sprintf("Showing %d Pokémon", len(pokemonList))))
	sb.WriteString("\n\n")

	// Show first 20 Pokemon
	displayCount := 20
	if len(pokemonList) < displayCount {
		displayCount = len(pokemonList)
	}

	for i := 0; i < displayCount; i++ {
		pokemon := pokemonList[i]
		// Format with proper spacing and styling
		line := fmt.Sprintf("  [%s] #%-3d %s",
			getCheckmark(i == m.selectedPokemon),
			pokemon.ID,
			capitalizeFirst(pokemon.Name))
		sb.WriteString(listItemStyle.Render(line))
		sb.WriteString("\n")
	}

	if len(pokemonList) > displayCount {
		sb.WriteString(listHelpStyle.Render(fmt.Sprintf("\nShowing 1-%d of %d", displayCount, len(pokemonList))))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(listHelpStyle.Render("↑/↓ or j/k: navigate • space/enter: select • 'l': lookup • 'q': quit"))

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

func getCheckmark(selected bool) string {
	if selected {
		return "x"
	}
	return " "
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
