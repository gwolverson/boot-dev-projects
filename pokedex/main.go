package main

import (
	"time"

	"github.com/gwolverson/pokedexcli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, time.Minute*5)
	config := Config{
		pokeapiClient: pokeClient,
		pokedex:       make(map[string]pokeapi.PokemonDetails),
	}
	startRepl(&config)
}
