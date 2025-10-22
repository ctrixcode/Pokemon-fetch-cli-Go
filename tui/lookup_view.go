package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrixcode/Pokemon-fetch-cli-Go/pokemon"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	lookupTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#FF6B6B")).
				Padding(0, 1).
				Bold(true)

	lookupInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	lookupHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	lookupErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true)

	lookupSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFA500")).
				Bold(true)

	lookupLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	// Title caser for proper capitalization
	titleCaser = cases.Title(language.English)
)

// LookupView represents the Pokemon lookup interface
type LookupView struct {
	input   textinput.Model
	result  *pokemon.PokemonData
	err     error
	loading bool
	width   int
	height  int
}

// pokemonFetchedMsg is sent when Pokemon data is fetched from the API
type pokemonFetchedMsg struct {
	pokemon *pokemon.PokemonData
	err     error
}

// NewLookupView creates a new lookup view
func NewLookupView() *LookupView {
	ti := textinput.New()
	ti.Placeholder = "Enter Pokémon name or ID (e.g., pikachu or 25)"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return &LookupView{
		input:  ti,
		width:  80,
		height: 24,
	}
}

// Init initializes the lookup view
func (v *LookupView) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the lookup view
func (v *LookupView) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = m.Width
		v.height = m.Height

	case tea.KeyMsg:
		switch m.String() {
		case "enter":
			if v.loading {
				return nil
			}

			query := strings.TrimSpace(v.input.Value())

			// If we have a result displayed, clear it and allow new search
			if v.result != nil {
				v.result = nil
				v.err = nil
				v.input.SetValue("")
				v.input.Placeholder = "Enter Pokémon name or ID (e.g., pikachu or 25)"
				v.input.Focus()
				return nil
			}

			// If no query, show error
			if query == "" {
				v.err = fmt.Errorf("please enter a Pokémon name or ID")
				return nil
			}

			// Fetch the pokemon
			v.loading = true
			v.err = nil
			return fetchPokemonCmd(query)

		case "esc":
			return func() tea.Msg { return switchToListViewMsg{} }

		case "ctrl+c":
			return tea.Quit
		}

	case pokemonFetchedMsg:
		v.loading = false
		v.err = m.err
		v.result = m.pokemon
		v.input.Placeholder = "Press Enter to search again or type a new name"
	}

	v.input, _ = v.input.Update(msg)
	return nil
}

// View renders the lookup view
func (v *LookupView) View() string {
	var s strings.Builder
	maxWidth := v.width
	if maxWidth < 60 {
		maxWidth = 60
	}
	if maxWidth > 100 {
		maxWidth = 100
	}

	// Title
	s.WriteString(lookupTitleStyle.Render("🔍 Pokémon Lookup"))
	s.WriteString("\n\n")

	// Loading
	if v.loading {
		s.WriteString("⏳ Fetching Pokémon data from API...\n\n")
		s.WriteString(lookupHelpStyle.Render("Press ESC to cancel"))
		return s.String()
	}

	// Error
	if v.err != nil {
		s.WriteString(lookupErrorStyle.Render(fmt.Sprintf("❌ Error: %v", v.err)))
		s.WriteString("\n\n")
	}

	// Input
	s.WriteString("Enter Pokémon name or ID:\n\n")
	s.WriteString(lookupInputStyle.Render(v.input.View()))
	s.WriteString("\n\n")

	// Result
	if v.result != nil {
		// Name + ID header - Make it more prominent
		divider := strings.Repeat("═", min(maxWidth-4, 60))
		s.WriteString(divider)
		s.WriteString("\n")
		s.WriteString(lookupSuccessStyle.Render(
			fmt.Sprintf("  %s (ID: #%d)", strings.ToUpper(v.result.Name), v.result.Id),
		))
		s.WriteString("\n")
		s.WriteString(divider)
		s.WriteString("\n\n")
		s.WriteString(v.formatPokemonData(maxWidth))
		s.WriteString("\n")
	}

	// Help text
	if v.result != nil {
		s.WriteString(lookupHelpStyle.Render("Press Enter to search again • ESC to go back"))
	} else {
		s.WriteString(lookupHelpStyle.Render("Press Enter to search • ESC to go back"))
	}

	return s.String()
}

// formatPokemonData formats the Pokemon data for display
func (v *LookupView) formatPokemonData(maxWidth int) string {
	if v.result == nil {
		return ""
	}

	p := v.result
	var s strings.Builder

	// Basic info
	s.WriteString(lookupLabelStyle.Render("📊 Basic Information:"))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("  Height: %d dm (%.1f m)\n", p.Height, float64(p.Height)/10))
	s.WriteString(fmt.Sprintf("  Weight: %d hg (%.1f kg)\n", p.Weight, float64(p.Weight)/10))
	s.WriteString(fmt.Sprintf("  Base Experience: %d\n\n", p.Base_experience))

	// Types
	if len(p.Types) > 0 {
		s.WriteString(lookupLabelStyle.Render("🏷️  Types:"))
		s.WriteString("\n")
		for _, t := range p.Types {
			s.WriteString(fmt.Sprintf("  • %s\n", titleCaser.String(t.Type.Name)))
		}
		s.WriteString("\n")
	}

	// Abilities
	if len(p.Abilities) > 0 {
		s.WriteString(lookupLabelStyle.Render("⚡ Abilities:"))
		s.WriteString("\n")
		for _, a := range p.Abilities {
			hidden := ""
			if a.IsHidden {
				hidden = " (Hidden)"
			}
			s.WriteString(fmt.Sprintf("  • %s%s\n", titleCaser.String(a.Ability.Name), hidden))
		}
		s.WriteString("\n")
	}

	// Stats
	if len(p.Stats) > 0 {
		s.WriteString(lookupLabelStyle.Render("📈 Base Stats:"))
		s.WriteString("\n")
		maxBarWidth := min(maxWidth-30, 40)
		for _, stat := range p.Stats {
			statName := titleCaser.String(stat.Stat.Name)
			barLength := (stat.Base_stat * maxBarWidth) / 255
			if barLength < 1 && stat.Base_stat > 0 {
				barLength = 1
			}
			if barLength > maxBarWidth {
				barLength = maxBarWidth
			}
			bar := strings.Repeat("█", barLength)
			s.WriteString(fmt.Sprintf("  %-20s %3d %s\n", statName+":", stat.Base_stat, bar))
		}
	}

	// Bottom divider
	s.WriteString("\n")
	s.WriteString(strings.Repeat("═", min(maxWidth-4, 60)))

	return s.String()
}

// fetchPokemonCmd creates a command to fetch Pokemon data from the API
func fetchPokemonCmd(nameOrID string) tea.Cmd {
	return func() tea.Msg {
		pokemonData, err := pokemon.FetchPokemon(nameOrID)
		if err != nil {
			return pokemonFetchedMsg{pokemon: nil, err: err}
		}
		return pokemonFetchedMsg{pokemon: pokemonData, err: nil}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
