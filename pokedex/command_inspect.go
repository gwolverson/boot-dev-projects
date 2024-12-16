package main

import (
	"fmt"
)

func inspect(config *Config, areaLocation string, pokemonName string) error {
	pokemon, exists := config.pokedex[pokemonName]

	if !exists {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")

		for _, stat := range pokemon.Stats {
			fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Printf("Types:\n")
		for _, types := range pokemon.Types {
			fmt.Printf("-%s\n", types.Type.Name)
		}
	}
	return nil
}
