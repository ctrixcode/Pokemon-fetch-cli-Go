package tui

import (
	"github.com/aviralgarg05/Pokemon-fetch-cli-Go/pokemon"
	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	mainMenuView viewState = iota
	listView
	detailView
)

type Model struct {
	state           viewState
	selectedPokemon int
	pokemonList     []int
	quitting        bool
}

func NewModel() Model {
	return Model{
		state:       listView,
		pokemonList: make([]int, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.state {
	case mainMenuView:
		return m.renderMainMenu()
	case listView:
		return m.renderListView()
	case detailView:
		return m.renderDetailView()
	}
	return ""
}

func (m Model) renderMainMenu() string {
	return "Pokemon Fetch CLI\n\nPress 'q' to quit\n"
}

func InitializeAndRun() error {
	pokemon.StoreData()
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
