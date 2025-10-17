package pokemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// PokemonData defines the structure of Pokémon API data
type PokemonData struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	Base_experience int `json:"base_experience"`
	Cries           struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height                   int    `json:"height"`
	Id                       int    `json:"id"`
	Is_default               bool   `json:"is_default"`
	Location_area_encounters string `json:"location_area_encounters"`
	Held_items               []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
	} `json:"held_items"`
	Name   string `json:"name"`
	Order  int    `json:"order"`
	Weight int    `json:"weight"`
	Moves  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		Back_default       string `json:"back_default"`
		Back_female        string `json:"back_female"`
		Back_shiny         string `json:"back_shiny"`
		Back_shiny_female  string `json:"back_shiny_female"`
		Front_default      string `json:"front_default"`
		Front_female       string `json:"front_female"`
		Front_shiny        string `json:"front_shiny"`
		Front_shiny_female string `json:"front_shiny_female"`
	} `json:"sprites"`
	Stats []struct {
		Base_stat int `json:"base_stat"`
		Effort    int `json:"effort"`
		Stat      struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

// Fetch all Pokémon data concurrently
func fetchAllPokemon() {
	const limit int = 1025
	var wq sync.WaitGroup
	wq.Add(limit)
	for index := 1; index <= limit; index++ {
		go fetchPokemonData(index, &wq)
	}
	wq.Wait()
}

// Fetch data for a single Pokémon and store it
func fetchPokemonData(id int, wq *sync.WaitGroup) {
	defer wq.Done()

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + fmt.Sprint(id))
	if err != nil {
		fmt.Println("error fetching Pokémon", id, ":", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, resp.Body) // drain body for connection reuse
		fmt.Println("unexpected status for Pokémon", id, ":", resp.Status)
		return
	}

	storePokemon(id, resp.Body)
}

// Store a Pokémon’s JSON data to disk
func storePokemon(id int, respBody io.ReadCloser) {
	body, err := io.ReadAll(respBody)
	if err != nil {
		fmt.Println("Error reading body for Pokémon", id, ":", err)
		return
	}

	var data PokemonData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling Pokémon JSON for ID", id, ":", err)
		return
	}

	dataFiltered, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error converting Pokémon data to JSON for ID", id, ":", err)
		return
	}

	// Get working directory for stable data path
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}
	dataDir := filepath.Join(cwd, "data")

	// Create data directory if missing
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return
	}

	filePath := filepath.Join(dataDir, fmt.Sprintf("%d.json", data.Id))
	if err := os.WriteFile(filePath, dataFiltered, 0o644); err != nil {
		fmt.Println("Error writing Pokémon JSON for ID", id, ":", err)
	}
}

// StoreData triggers fetching all Pokémon if not already initialized
func StoreData() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	dataDir := filepath.Join(cwd, "data")
	initMarker := filepath.Join(dataDir, "init.txt")

	if _, err := os.Stat(initMarker); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Fetching Pokémon data for the first time... this may take a while.")
		fetchAllPokemon()

		if err := os.WriteFile(initMarker, []byte("initialized"), 0o644); err != nil {
			fmt.Println("Error creating init marker file:", err)
		} else {
			fmt.Println("All Pokémon data stored successfully.")
		}
	} else {
		fmt.Println("Pokémon data already initialized.")
	}
}
