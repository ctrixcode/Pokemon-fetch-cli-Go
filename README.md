# Pokémon Fetch CLI (Go)

A simple Go command-line tool (`pokemon-fetch-cli-go`) to fetch Pokémon data from the [PokéAPI](https://pokeapi.co/) and store it locally.

## How It Works

When you run the application, it fetches data for the first 151 Pokémon from the PokéAPI and saves it into a `pokedex.json` file in the project's root directory.

## Getting Started

To run the application, follow these steps:

1.  **Clone the repository**
    ```sh
    git clone <repository-url>
    ```

2.  **Run the application**
    ```sh
    go run main.go
    ```
    After running, you will find a `pokedex.json` file in the project directory containing the fetched data.

## Contributing

This project is just getting started! The goal is to add more features for processing and exporting the stored Pokémon data.

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) and check out the [TODO list](./docs/todo.md) for planned features.
