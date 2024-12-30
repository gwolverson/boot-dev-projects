package main

import (
	"github.com/gwolverson/go-courses/blog-aggregator/internal/config"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/database"
	"log"
)

type state struct {
	appConfig *config.AppConfig
	db        *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (commands *commands) register(name string, function func(*state, command) error) {
	commands.commandMap[name] = function
}

func (commands *commands) run(state *state, command command) error {
	cmdFunc, exists := commands.commandMap[command.name]
	if !exists {
		log.Fatalf("Command not found: %s\n", command.name)
	}

	return cmdFunc(state, command)
}
