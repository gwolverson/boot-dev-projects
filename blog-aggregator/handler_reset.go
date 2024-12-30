package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(state *state, cmd command) error {
	err := state.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Printf("Error deleting users: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully deleted users")
	return nil
}
