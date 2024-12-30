package main

import (
	"database/sql"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/config"
	"github.com/gwolverson/go-courses/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatalf("At least two arguments are required; the command and a parameter")
	}

	cmd := command{
		name: userArgs[1],
		args: userArgs[2:],
	}

	db, err := sql.Open("postgres", appConfig.DbUrl)
	if err != nil {
		log.Fatalf("Error opening database connection: %v\n", err)
	}
	dbQueries := database.New(db)

	appState := state{
		appConfig: &appConfig,
		db:        dbQueries,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("register", handlerRegister)
	cmds.register("login", handlerLogin)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerRSS)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", handlerFollowing)
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	err = cmds.run(&appState, cmd)
	if err != nil {
		log.Fatalf("Error running command: %v\n", err)
	}
}
