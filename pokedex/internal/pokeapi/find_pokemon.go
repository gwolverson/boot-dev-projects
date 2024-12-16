package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type LocationPokemon struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string
	Url  string
}

func (c *Client) FindPokemon(name string) ([]string, error) {
	url := baseURL + "/location-area/" + name

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var locationPokemon LocationPokemon
	err = json.Unmarshal(dat, &locationPokemon)
	if err != nil {
		return []string{}, err
	}

	var pokemon []string
	for _, encounter := range locationPokemon.PokemonEncounters {
		pokemon = append(pokemon, encounter.Pokemon.Name)
	}

	return pokemon, nil
}
