package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type PokemonDetails struct {
	BaseExperience int `json:"base_experience"`
	Name           string
	Height         int
	Weight         int
	Stats          []BaseStat
	Types          []PokemonType
}

type BaseStat struct {
	BaseStat int `json:"base_stat"`
	Effort   int
	Stat     Stat
}

type Stat struct {
	Name string
}

type PokemonType struct {
	Slot int
	Type Type
}

type Type struct {
	Name string
}

func (c *Client) GetPokemon(name string) (PokemonDetails, error) {
	url := baseURL + "/pokemon/" + name

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonDetails{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDetails{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDetails{}, err
	}

	var pokemon PokemonDetails
	err = json.Unmarshal(dat, &pokemon)
	if err != nil {
		return PokemonDetails{}, err
	}

	return pokemon, nil
}
