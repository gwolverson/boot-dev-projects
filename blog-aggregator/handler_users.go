package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func handlerLogin(state *state, command command) error {
	if len(command.args) == 0 {
		return errors.New("login command expects a single <username> arg")
	}

	_, err := state.db.GetUser(context.Background(), command.args[0])
	if err != nil {
		fmt.Printf("user %s does not exist\n", command.args[0])
		os.Exit(1)
	}

	err = state.appConfig.SetUser(command.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Login command successfully executed: %s\n", command.args[0])
	return nil
}

func handlerGetUsers(state *state, command command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error encountered when retrieving users")
		return err
	}

	currentUsername := state.appConfig.CurrentUserName
	for _, user := range users {
		if currentUsername == user.Name {
			fmt.Printf("* %s (current) \n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
