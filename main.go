package main

import (
	"fmt"
	"os"

	"github.com/faxter/bloggator/internal/command"
	"github.com/faxter/bloggator/internal/config"
	"github.com/faxter/bloggator/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	st := state.NewState(&cfg)
	cmds := command.NewCommandSet()
	cmds.RegisterBuiltIns()

	if len(os.Args) < 2 {
		fmt.Println("missing a command!")
		os.Exit(1)
	}

	command := command.NewCommand(os.Args[1])
	if len(os.Args) > 1 {
		command.Args = os.Args[2:]
	}

	err = cmds.Run(st, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
