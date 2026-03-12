package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/faxter/bloggator/internal/command"
	"github.com/faxter/bloggator/internal/config"
	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/state"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)

	st := state.NewState(dbQueries, &cfg)
	commandSet := command.NewCommandSet()
	commandSet.RegisterBuiltIns()

	if len(os.Args) < 2 {
		fmt.Println("missing a command!")
		os.Exit(1)
	}

	command := command.NewCommand(os.Args[1])
	if len(os.Args) > 1 {
		command.Args = os.Args[2:]
	}

	err = commandSet.Run(st, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
