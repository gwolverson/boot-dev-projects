package main

import (
	"fmt"
	"os"
)

func commandExit(config *Config, areaLocation string, pokemonName string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
