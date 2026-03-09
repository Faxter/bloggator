package main

import (
	"fmt"
	"os"

	"github.com/faxter/bloggator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, found := c.commands[cmd.name]
	if !found {
		return fmt.Errorf("could not find command %s", cmd.name)
	}
	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login requires argument for username")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been set to", cmd.args[0])
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	st := &state{config: &cfg}
	cmds := commands{commands: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("missing a command!")
		os.Exit(1)
	}

	command := command{name: os.Args[1]}
	if len(os.Args) > 1 {
		command.args = os.Args[2:]
	}

	err = cmds.run(st, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
