package pokemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

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
		Url  string `json:"url"`
	} `json:"forms"`
	Height                   int    `json:"height"`
	Id                       int    `json:"id"`
	Is_default               bool   `json:"is_default"`
	Location_area_encounters string `json:"location_area_encounters"`
	Held_items               []struct {
		Item struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"item"`
	} `json:"held_items"`
	Name   string `json:"name"`
	Order  int    `json:"order"`
	Weight int    `json:"weight"`
	Moves  []struct {
		Move struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		Url  string `json:"url"`
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
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

// FetchPokemon fetches a single Pokemon by name or ID from the PokeAPI
func FetchPokemon(nameOrID string) (*PokemonData, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", nameOrID)

	resp, err := http.Get(url)
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
	err = json.Unmarshal(body, &data)
	if err != nil {
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

func fetchPokemonData(id int, wq *sync.WaitGroup) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + fmt.Sprint(id))
	if err != nil {
		fmt.Println("error with ", id, ": ", err)
	}
	storePokemon(resp.Body)
	wq.Done()
}

func storePokemon(respBody io.ReadCloser) {
	body, err := io.ReadAll(respBody)
	if err != nil {
		fmt.Println("Error during Converting body into byte", ": ", err)
		return
	}
	var data PokemonData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	dataFiltered, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error trying to convert filtered data into byte", data.Id, ": ", err)
		return
	}

	err = os.WriteFile("data/"+fmt.Sprint(data.Id)+".json", dataFiltered, 0644)
	if err != nil {
		fmt.Println("Error writing file for Pokemon", data.Id, ":", err)
	}
}

func StoreData() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Println("Error creating data directory:", err)
		return
	}

	// Create pokemon directory if it doesn't exist
	if err := os.MkdirAll("pokemon", 0755); err != nil {
		fmt.Println("Error creating pokemon directory:", err)
		return
	}

	// Check if we've already fetched data
	if _, err := os.Stat("pokemon/init.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Fetching all Pokémon data (this may take a while)...")
		fetchAllPokemon()
		fmt.Println("Fetch complete! Creating init marker file...")

		_, err := os.Create("pokemon/init.txt")
		if err != nil {
			fmt.Println("Warning: file cannot be created:", err)
		}
	} else {
		fmt.Println("Data already fetched (found pokemon/init.txt)")
	}
}
