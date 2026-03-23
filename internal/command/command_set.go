package command

import (
	"fmt"

	"github.com/faxter/bloggator/internal/state"
)

type CommandSet struct {
	Commands map[string]func(*state.State, Command) error
}

func NewCommandSet() CommandSet {
	return CommandSet{Commands: make(map[string]func(*state.State, Command) error)}
}

func (c *CommandSet) RegisterBuiltIns() {
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregate)
	c.register("addfeed", withUserLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", withUserLoggedIn(handlerFollow))
	c.register("following", withUserLoggedIn(handlerFollows))
	c.register("unfollow", withUserLoggedIn(handlerUnfollow))
	c.register("browse", withUserLoggedIn(handlerBrowse))
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
