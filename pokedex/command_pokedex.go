package main

import (
	"fmt"
)

func pokedex(config *Config, areaLocation string, pokemonName string) error {
	fmt.Println("Your Pokedex:")
	for key := range config.pokedex {
		fmt.Printf("- %s\n", key)
	}
	return nil
}
