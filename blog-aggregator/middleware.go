package main

import (
	"context"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(state *state, cmd command, user database.User) error) func(*state, command) error {
	return func(state *state, cmd command) error {
		user, err := state.db.GetUser(context.Background(), state.appConfig.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(state, cmd, user)
	}
}
