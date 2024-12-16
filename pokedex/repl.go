package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gwolverson/pokedexcli/internal/pokeapi"
)

type Config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	pokedex          map[string]pokeapi.PokemonDetails
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, areaLocation string, pokemonName string) error
}

func startRepl(config *Config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := getCommands(config)[commandName]

		var location string
		if commandName == "explore" {
			location = words[1]
		}

		var pokemonName string
		if commandName == "catch" || commandName == "inspect" {
			pokemonName = words[1]
		}

		if exists {
			err := command.callback(config, location, pokemonName)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands(config *Config) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display names of 20 location areas in the Pokemon world",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Display names of 20 previous location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Display all the pokemon located in the given area",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a specific pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon you have already caught",
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all pokemon in your pokedex",
			callback:    pokedex,
		},
	}
}
