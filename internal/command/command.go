package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/state"
	"github.com/google/uuid"
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
	c.register("register", handlerRegister)
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
	username := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user has been set to", username)
	return nil
}

func handlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login requires argument for username")
	}
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]})
	if err != nil {
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("user was created:", user)
	return nil
}
