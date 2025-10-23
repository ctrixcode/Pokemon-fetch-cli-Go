package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// viewMode represents the current view
type viewMode int

const (
	listViewMode viewMode = iota
	lookupViewMode
)

// PokemonListItem represents a single Pokemon in the list
type PokemonListItem struct {
	ID   string
	Name string
	File string
}

// Model represents the TUI state for the Pokemon list
type Model struct {
	pokemonList []PokemonListItem
	cursor      int
	selected    map[int]struct{}
	width       int
	height      int
	loading     bool
	err         error
	currentView viewMode
	lookupView  *LookupView
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
		currentView: listViewMode,
	}
}

// Init is the first function that will be called when the program starts
func (m Model) Init() tea.Cmd {
	return LoadPokemonFromFiles
}

// Update handles user input and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case switchToListViewMsg:
		m.currentView = listViewMode
		return m, nil

	case switchToLookupViewMsg:
		// Initialize lookup view if it doesn't exist
		if m.lookupView == nil {
			m.lookupView = NewLookupView()
		} else {
			// Reset the lookup view to clear any previous state
			m.lookupView.Reset()
		}
		m.currentView = lookupViewMode
		cmd := m.lookupView.Init()
		return m, cmd
	}

	// Delegate to lookup view if active
	if m.currentView == lookupViewMode && m.lookupView != nil {
		cmd := m.lookupView.Update(msg)
		return m, cmd
	}

	// Handle list view updates
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

		case "l", "L":
			// Switch to lookup view with message dispatch
			return m, func() tea.Msg { return switchToLookupViewMsg{} }

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
				if _, ok := m.selected[m.cursor]; ok {
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
	if m.currentView == lookupViewMode && m.lookupView != nil {
		return m.lookupView.View()
	}

	if m.loading {
		return "\n Loading Pokémon data...\n\n"
	}

	if m.err != nil {
		return fmt.Sprintf("\n Error loading Pokémon data: %v\n\n", m.err)
	}

	if len(m.pokemonList) == 0 {
		return "\n No Pokémon data found. Please run the data fetch first.\n\n"
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("Pokémon List"))
	s.WriteString("\n\n")

	start := 0
	end := len(m.pokemonList)
	maxVisible := m.height - 8

	if maxVisible > 0 && len(m.pokemonList) > maxVisible {
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
			s.WriteString(selectedItemStyle.Render(line))
		} else {
			s.WriteString(itemStyle.Render(line))
		}
		s.WriteString("\n")
	}

	if maxVisible > 0 && len(m.pokemonList) > maxVisible {
		s.WriteString(paginationStyle.Render(fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(m.pokemonList))))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓ or j/k: navigate • space/enter: select • L: lookup • q: quit • home/end: first/last • pgup/pgdn: page"))
	s.WriteString("\n")

	return s.String()
}

// pokemonDataMsg is used to pass loaded Pokemon data to the model
type pokemonDataMsg struct {
	pokemonList []PokemonListItem
	err         error
}

// Message types for view switching
type switchToListViewMsg struct{}
type switchToLookupViewMsg struct{}
