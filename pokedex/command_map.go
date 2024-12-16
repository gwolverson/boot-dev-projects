package main

import (
	"errors"
	"fmt"
)

const baseUrl = "https://pokeapi.co/api/v2/location-area"

type Location struct {
	Count    int
	Next     string
	Previous string
	Results  []LocationResult
}

type LocationResult struct {
	Name string
	Url  string
}

func commandMapf(cfg *Config, areaLocation string, pokemonName string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *Config, areaLocation string, pokemonName string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
