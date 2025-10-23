package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	input    textinput.Model
	viewport viewport.Model
	result   *pokemon.PokemonData
	err      error
	loading  bool
	width    int
	height   int
	ready    bool
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

	vp := viewport.New(80, 16)
	vp.YPosition = 8

	return &LookupView{
		input:    ti,
		viewport: vp,
		width:    80,
		height:   24,
		ready:    true,
	}
}

func (v *LookupView) Init() tea.Cmd {
	return textinput.Blink
}

func (v *LookupView) Reset() {
	v.result = nil
	v.err = nil
	v.loading = false
	v.input.SetValue("")
	v.viewport.SetContent("")
	v.viewport.GotoTop()
	v.input.Focus()
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
	var cmds []tea.Cmd

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = m.Width
		v.height = m.Height

		// Adjust input width based on window size (clamped 20..70)
		inputWidth := clamp(m.Width-10, 20, 70)
		v.input.Width = inputWidth

		// Adjust viewport size (clamped)
		headerHeight := 8 // Approximate height for title, input, and help text
		v.viewport.Width = clamp(m.Width, 20, 200)
		vh := m.Height - headerHeight
		if vh < 3 {
			vh = 3
		}
		v.viewport.Height = vh
		v.viewport.YPosition = headerHeight

		if v.result != nil {
			v.viewport.SetContent(v.getResultContent())
		}
		return nil

	case tea.KeyMsg:
		switch m.String() {
		case "enter":
			if v.result != nil {
				v.result = nil
				v.err = nil
				v.viewport.SetContent("")
				v.viewport.GotoTop()
				v.input.SetValue("")
				v.input.Focus()
				return textinput.Blink
			}
			if strings.TrimSpace(v.input.Value()) != "" {
				v.loading = true
				v.err = nil
				v.result = nil
				v.viewport.SetContent("")
				v.viewport.GotoTop()
				return v.fetchPokemonData()
			}
			return nil
		case "esc", "ctrl+c":
			v.Reset()
			return func() tea.Msg { return switchToListViewMsg{} }
		case "up", "down", "pgup", "pgdown", "home", "end":
			if v.result != nil {
				v.viewport, cmd = v.viewport.Update(m)
				return cmd
			}
			return nil
		}
		v.input, cmd = v.input.Update(m)
		return cmd

	case pokemonFetchedMsg:
		v.loading = false
		v.err = m.err
		v.result = m.pokemon

		v.viewport.SetContent(v.getResultContent())
		v.viewport.GotoTop()
		v.input.SetValue("")
		v.input.Focus()
		return textinput.Blink

	default:
		if v.result != nil {
			v.viewport, cmd = v.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
		v.input, cmd = v.input.Update(msg)
		cmds = append(cmds, cmd)
		return tea.Batch(cmds...)
	}
}

func (v *LookupView) getResultContent() string {
	if v.result == nil {
		return ""
	}

	var s strings.Builder
	maxWidth := v.width
	if maxWidth < 60 {
		maxWidth = 60
	}
	if maxWidth > 120 {
		maxWidth = 120
	}

	contentWidth := min(maxWidth-4, 80)
	divider := strings.Repeat("═", contentWidth)

	s.WriteString("\n")
	s.WriteString(divider + "\n")
	s.WriteString(lookupSuccessStyle.Render(
		fmt.Sprintf("  %s (ID: #%d)", strings.ToUpper(v.result.Name), v.result.Id),
	))
	s.WriteString("\n")
	s.WriteString(divider + "\n\n")

	s.WriteString(v.formatPokemonData(maxWidth))
	s.WriteString("\n")
	s.WriteString(divider + "\n")

	return s.String()
}

func (v *LookupView) View() string {
	var s strings.Builder

	title := lookupTitleStyle.Render("🔍 Pokémon Lookup")
	s.WriteString(title + "\n\n")

	if v.loading {
		s.WriteString("⏳ Fetching Pokémon data from API...\n\n")
		s.WriteString(lookupHelpStyle.Render("Press ESC to cancel"))
		return s.String()
	}

	if v.err != nil {
		s.WriteString(lookupErrorStyle.Render(fmt.Sprintf("❌ Error: %v", v.err)))
		s.WriteString("\n\n")
	}

	if v.result == nil {
		s.WriteString("Enter Pokémon name or ID:\n\n")
		s.WriteString(lookupInputStyle.Render(v.input.View()))
		s.WriteString("\n\n")
		s.WriteString(lookupHelpStyle.Render("Press Enter to search • ESC to go back"))
	} else {
		s.WriteString(v.viewport.View())
		s.WriteString("\n")
		s.WriteString(lookupHelpStyle.Render("↑/↓: Scroll • Enter: New search • ESC: Go back"))
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
		if maxBarWidth < 10 {
			maxBarWidth = 10
		}
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

	return s.String()
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}
