package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const quitMessage = "\nPress 'q' to quit\n"

type PokemonDetail struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Types          []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
	} `json:"abilities"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func (m Model) renderDetailView() string {
	var sb strings.Builder

	sb.WriteString("Pokemon Detail View\n")
	sb.WriteString("===================\n\n")

	if m.selectedPokemon <= 0 {
		sb.WriteString("No Pokemon selected.\n")
		sb.WriteString(quitMessage)
		return sb.String()
	}

	pokemon := m.getPokemonDetail(m.selectedPokemon)

	if pokemon == nil {
		sb.WriteString(fmt.Sprintf("Pokemon with ID %d not found.\n", m.selectedPokemon))
		sb.WriteString(quitMessage)
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Name: %s\n", strings.Title(pokemon.Name)))
	sb.WriteString(fmt.Sprintf("ID: %d\n", pokemon.ID))
	sb.WriteString(fmt.Sprintf("Height: %d\n", pokemon.Height))
	sb.WriteString(fmt.Sprintf("Weight: %d\n", pokemon.Weight))
	sb.WriteString(fmt.Sprintf("Base Experience: %d\n\n", pokemon.BaseExperience))

	if len(pokemon.Types) > 0 {
		sb.WriteString("Types:\n")
		for _, t := range pokemon.Types {
			sb.WriteString(fmt.Sprintf("  - %s\n", strings.Title(t.Type.Name)))
		}
		sb.WriteString("\n")
	}

	if len(pokemon.Abilities) > 0 {
		sb.WriteString("Abilities:\n")
		for _, a := range pokemon.Abilities {
			hidden := ""
			if a.IsHidden {
				hidden = " (Hidden)"
			}
			sb.WriteString(fmt.Sprintf("  - %s%s\n", strings.Title(a.Ability.Name), hidden))
		}
		sb.WriteString("\n")
	}

	if len(pokemon.Stats) > 0 {
		sb.WriteString("Stats:\n")
		for _, s := range pokemon.Stats {
			sb.WriteString(fmt.Sprintf("  - %s: %d\n", strings.Title(s.Stat.Name), s.BaseStat))
		}
	}

	sb.WriteString(quitMessage)

	return sb.String()
}

func (m Model) getPokemonDetail(id int) *PokemonDetail {
	dataPath := filepath.Join("data", fmt.Sprintf("%d.json", id))

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil
	}

	var pokemon PokemonDetail
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return nil
	}

	return &pokemon
}
