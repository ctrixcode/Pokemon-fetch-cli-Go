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
		// seems useless
		// Version_details []struct {
		// 	Rarity  int `json:"rarity"`
		// 	Version struct {
		// 		Name string `json:"name"`
		// 		Url  string `json:"url"`
		// 	}
		// } `json:"version_details"`
	} `json:"held_items"`
	Name   string `json:"name"`
	Order  int    `json:"order"`
	Weight int    `json:"weight"`
	Moves  []struct {
		Move struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"move"`
		// useful to know moves learnt first generation
		// Version_group_details []struct {
		// 	Level_learned_at  int `json:"level_learned_at"`
		// 	Move_learn_method struct {
		// 		Name string `json:"name"`
		// 		Url  string `json:"url"`
		// 	} `json:"move_learn_method"`
		// 	Version_group struct {
		// 		Name string `json:"name"`
		// 		Url  string `json:"url"`
		// 	} `json:"version_group"`
		// } `json:"version_group_details"`
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
	}
	var data PokemonData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("err:", err)
	}
	dataFiltered, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error trying to convert filtered data into byte", data.Id, ": ", err)
	}
	// fmt.Println(data.Id, data.Name)
	os.WriteFile("data/"+fmt.Sprint(data.Id)+".json", dataFiltered, 0655)
}
func StoreData() {
	if _, err := os.Stat("pokemon/init.txt"); errors.Is(err, os.ErrNotExist) {
		fetchAllPokemon()
		_, err := os.Create("pokemon/init.txt")
		if err != nil {
			println("file cannot be created: ", err)
		}
	}
}
