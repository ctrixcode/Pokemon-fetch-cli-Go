package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PokemonListItem represents a single Pokemon in the list
type PokemonListItem struct {
	ID   string
	Name string
	File string
}

// Model represents the TUI state for the Pokemon list
type Model struct {
	pokemonList  []PokemonListItem
	cursor       int
	selected     map[int]struct{}
	width        int
	height       int
	loading      bool
	err          error
}

// Styling
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	paginationStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingRight(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// NewModel creates a new model for the Pokemon list TUI
func NewModel() Model {
	return Model{
		pokemonList: []PokemonListItem{},
		selected:    make(map[int]struct{}),
		loading:     true,
	}
}

// Init is the first function that will be called when the program starts
func (m Model) Init() tea.Cmd {
	return LoadPokemonFromFiles
}

// Update handles user input and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case pokemonDataMsg:
		m.pokemonList = msg.pokemonList
		m.loading = false
		m.err = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.pokemonList)-1 {
				m.cursor++
			}

		case "enter", " ":
			if len(m.pokemonList) > 0 {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		case "home":
			m.cursor = 0

		case "end":
			if len(m.pokemonList) > 0 {
				m.cursor = len(m.pokemonList) - 1
			}

		case "pgup":
			m.cursor -= 10
			if m.cursor < 0 {
				m.cursor = 0
			}

		case "pgdown":
			m.cursor += 10
			if m.cursor >= len(m.pokemonList) {
				m.cursor = len(m.pokemonList) - 1
			}
		}
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.loading {
		return "\n Loading Pokémon data...\n\n"
	}

	if m.err != nil {
		return fmt.Sprintf("\n Error loading Pokémon data: %v\n\n", m.err)
	}

	if len(m.pokemonList) == 0 {
		return "\n No Pokémon data found. Please run the data fetch first.\n\n"
	}

	// Header
	s := titleStyle.Render("Pokémon List") + "\n\n"

	// Calculate visible range for pagination
	start := 0
	end := len(m.pokemonList)
	maxVisible := m.height - 8 // Reserve space for header, help, and margins

	if maxVisible > 0 && len(m.pokemonList) > maxVisible {
		// Center the cursor in the visible area
		start = m.cursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(m.pokemonList) {
			end = len(m.pokemonList)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}
	}

	// Render Pokemon list
	for i := start; i < end; i++ {
		pokemon := m.pokemonList[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✓"
		}

		line := fmt.Sprintf("%s [%s] #%s %s", cursor, checked, pokemon.ID, strings.Title(pokemon.Name))

		if m.cursor == i {
			s += selectedItemStyle.Render(line) + "\n"
		} else {
			s += itemStyle.Render(line) + "\n"
		}
	}

	// Pagination info
	if maxVisible > 0 && len(m.pokemonList) > maxVisible {
		s += paginationStyle.Render(fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(m.pokemonList))) + "\n"
	}

	// Help
	s += "\n" + helpStyle.Render("↑/↓ or j/k: navigate • space/enter: select • q: quit • home/end: first/last • pgup/pgdn: page") + "\n"

	return s
}

// pokemonDataMsg is used to pass loaded Pokemon data to the model
type pokemonDataMsg struct {
	pokemonList []PokemonListItem
	err         error
}

