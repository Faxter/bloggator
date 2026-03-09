package command

import (
	"fmt"

	"github.com/faxter/bloggator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type CommandSet struct {
	Commands map[string]func(*state.State, Command) error
}

func NewCommand(name string) Command {
	return Command{Name: name}
}

func NewCommandSet() CommandSet {
	return CommandSet{Commands: make(map[string]func(*state.State, Command) error)}
}

func (c *CommandSet) RegisterBuiltIns() {
	c.register("login", handlerLogin)
}

func (c *CommandSet) Run(s *state.State, cmd Command) error {
	command, found := c.Commands[cmd.Name]
	if !found {
		return fmt.Errorf("could not find command %s", cmd.Name)
	}
	return command(s, cmd)
}

func (c *CommandSet) register(name string, f func(*state.State, Command) error) {
	c.Commands[name] = f
}

func handlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login requires argument for username")
	}
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been set to", cmd.Args[0])
	return nil
}
