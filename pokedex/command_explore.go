package main

import "fmt"

func explore(config *Config, areaLocation string, pokemonName string) error {
	pokemonResp, err := config.pokeapiClient.FindPokemon(areaLocation)

	if err != nil {
		return err
	}

	for _, pokemon := range pokemonResp {
		fmt.Printf("%s\n", pokemon)
	}
	return nil
}
