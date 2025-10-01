# Pokémon Fetch CLI (Go)

A Go command-line tool (`pokemon-fetch-cli-go`) to fetch Pokémon data from the [PokéAPI](https://pokeapi.co/) and display it in an interactive Terminal User Interface (TUI).

## Features

- 🎯 **Data Fetching**: Downloads data for all 1025+ Pokémon from the PokéAPI
- 📋 **Interactive List**: Browse Pokémon in a scrollable, navigable list within the terminal
- ⌨️ **Keyboard Navigation**: Use arrow keys, vim keys (j/k), or page up/down to navigate
- 🔍 **Selection**: Mark/unmark Pokémon with space or enter keys
- 💾 **Local Storage**: Data is stored in individual JSON files for each Pokémon

## How It Works

The application fetches data for all available Pokémon from the PokéAPI and stores each Pokémon's data in separate JSON files in a `data/` directory. You can then use the interactive TUI to browse through the collected Pokémon data.

## Getting Started

### Prerequisites
- Go 1.22.2 or later

### Installation

1. **Clone the repository**
   ```sh
   git clone <repository-url>
   cd pokemon-fetch-cli-go
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

### Usage

#### Fetch Pokémon Data
```sh
go run main.go --fetch
```
This downloads all Pokémon data and stores it in the `data/` directory.

#### View Interactive List
```sh
go run main.go --list
```
This opens the interactive TUI to browse through the Pokémon data.

#### Fetch and View (Combined)
```sh
go run main.go --fetch --list
```
This fetches the data first, then immediately opens the interactive list.

#### Default Behavior (Backward Compatibility)
```sh
go run main.go
```
Without flags, the application will fetch data (maintaining backward compatibility).

## Interactive List Controls

- **↑/↓** or **j/k**: Navigate up/down
- **Space/Enter**: Select/deselect current Pokémon
- **Home/End**: Jump to first/last Pokémon
- **Page Up/Page Down**: Navigate by pages
- **q** or **Ctrl+C**: Quit the application

## Contributing

This project is just getting started! The goal is to add more features for processing and exporting the stored Pokémon data.

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) and check out the [TODO list](./docs/todo.md) for planned features.
