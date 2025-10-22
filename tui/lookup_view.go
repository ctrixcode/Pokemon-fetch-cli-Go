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

	titleCaser = cases.Title(language.English)
)

type LookupView struct {
	input   textinput.Model
	result  *pokemon.PokemonData
	err     error
	loading bool
	width   int
	height  int
}

type pokemonFetchedMsg struct {
	pokemon *pokemon.PokemonData
	err     error
}

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

func (v *LookupView) Init() tea.Cmd {
	return textinput.Blink
}

func (v *LookupView) fetchPokemonData() tea.Cmd {
	input := strings.TrimSpace(v.input.Value())
	if input == "" {
		return nil
	}

	return func() tea.Msg {
		data, err := pokemon.FetchPokemon(input)
		return pokemonFetchedMsg{pokemon: data, err: err}
	}
}

func (v *LookupView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m := msg.(type) {
	case tea.KeyMsg:
		// Always update the input model first for proper async behavior
		var inputCmd tea.Cmd
		v.input, inputCmd = v.input.Update(m)
		cmd = inputCmd

		switch m.String() {
		case "enter":
			return tea.Batch(cmd, v.fetchPokemonData())
		case "esc", "ctrl+c":
			return tea.Batch(cmd, func() tea.Msg { return switchToListViewMsg{} })
		default:
			return cmd
		}

	default:
		v.input, cmd = v.input.Update(msg)
		return cmd
	}
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

	if v.loading {
		s.WriteString("⏳ Fetching Pokémon data from API...\n\n")
		s.WriteString(lookupHelpStyle.Render("Press ESC to cancel"))
		return s.String()
	}

	if v.err != nil {
		s.WriteString(lookupErrorStyle.Render(fmt.Sprintf("❌ Error: %v", v.err)))
		s.WriteString("\n\n")
	}

	s.WriteString("Enter Pokémon name or ID:\n\n")
	s.WriteString(lookupInputStyle.Render(v.input.View()))
	s.WriteString("\n\n")

	if v.result != nil {
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

	if v.result != nil {
		s.WriteString(lookupHelpStyle.Render("Press Enter to search again • ESC to go back"))
	} else {
		s.WriteString(lookupHelpStyle.Render("Press Enter to search • ESC to go back"))
	}

	return s.String()
}

func (v *LookupView) formatPokemonData(maxWidth int) string {
	if v.result == nil {
		return ""
	}

	p := v.result
	var s strings.Builder

	s.WriteString(lookupLabelStyle.Render("📊 Basic Information:"))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("  Height: %d dm (%.1f m)\n", p.Height, float64(p.Height)/10))
	s.WriteString(fmt.Sprintf("  Weight: %d hg (%.1f kg)\n", p.Weight, float64(p.Weight)/10))
	s.WriteString(fmt.Sprintf("  Base Experience: %d\n\n", p.Base_experience))

	if len(p.Types) > 0 {
		s.WriteString(lookupLabelStyle.Render("🏷️  Types:"))
		s.WriteString("\n")
		for _, t := range p.Types {
			s.WriteString(fmt.Sprintf("  • %s\n", titleCaser.String(t.Type.Name)))
		}
		s.WriteString("\n")
	}

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

	s.WriteString("\n")
	s.WriteString(strings.Repeat("═", min(maxWidth-4, 60)))
	return s.String()
}

func fetchPokemonCmd(nameOrID string) tea.Cmd {
	return func() tea.Msg {
		pokemonData, err := pokemon.FetchPokemon(nameOrID)
		if err != nil {
			return pokemonFetchedMsg{pokemon: nil, err: err}
		}
		return pokemonFetchedMsg{pokemon: pokemonData, err: nil}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
