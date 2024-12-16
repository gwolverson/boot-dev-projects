package main

import (
	"fmt"
	"math/rand"
)

func catch(config *Config, locationArea string, pokemonName string) error {
	pokemonResp, err := config.pokeapiClient.GetPokemon(pokemonName)

	if err != nil {
		return err
	}

	chance := rand.Intn(int(pokemonResp.BaseExperience))

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if chance > 40 {
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	config.pokedex[pokemonName] = pokemonResp

	return nil
}
