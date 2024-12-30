package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/database"
	"github.com/lib/pq"
	"log"
	"os"
	"time"
)

func handlerRegister(state *state, command command) error {
	if len(command.args) == 0 {
		return errors.New("register command expects a single <username> arg")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      command.args[0],
	}

	registeredUser, err := state.db.CreateUser(context.Background(), newUser)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Username already exists.")
			os.Exit(1)
		} else {
			fmt.Printf("Failed to register user: %v\n", err)
			os.Exit(1)
		}
	}
	err = state.appConfig.SetUser(registeredUser.Name)
	if err != nil {
		log.Fatalf("Failed to update current user: %v", err)
		return err
	}

	fmt.Printf("User registered successfully! %v\n", registeredUser)
	return nil
}
