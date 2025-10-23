package pokemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
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

// FetchPokemon fetches a single Pokemon by name or ID from the PokeAPI
func FetchPokemon(nameOrID string) (*PokemonData, error) {
	safeName := url.PathEscape(nameOrID)
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", safeName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokemon not found: %s (status: %d)", nameOrID, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data PokemonData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse pokemon data: %w", err)
	}

	return &data, nil
}

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

	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error fetching", id, ":", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("pokemon", id, "fetch failed with status:", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body for", id, ":", err)
		return
	}

	storePokemon(body)
}

func storePokemon(body []byte) {
	var data PokemonData
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("err:", err)
		return
	}

	dataFiltered, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error trying to convert filtered data into byte", data.Id, ":", err)
		return
	}

	if err := os.WriteFile("data/"+fmt.Sprint(data.Id)+".json", dataFiltered, 0644); err != nil {
		fmt.Println("Error writing file for Pokemon", data.Id, ":", err)
	}
}

func StoreData() {
	// Create directories if they don't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return
	}
	if err := os.MkdirAll("pokemon", 0755); err != nil {
		fmt.Println("Error creating pokemon directory:", err)
		return
	}

	// Check if data already fetched
	if _, err := os.Stat("pokemon/init.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Fetching all Pokémon data (this may take a while)...")
		fetchAllPokemon()
		fmt.Println("Fetch complete! Creating init marker file...")

		if err := os.WriteFile("pokemon/init.txt", []byte{}, 0644); err != nil {
			fmt.Println("Warning: could not create init marker:", err)
		}
	} else {
		fmt.Println("Data already fetched (found pokemon/init.txt)")
	}
}
