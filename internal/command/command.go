package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/rss"
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
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregate)
	c.register("addfeed", withUserLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", withUserLoggedIn(handlerFollow))
	c.register("following", withUserLoggedIn(handlerFollows))
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
	fmt.Println("user was created:", user.Name)
	return nil
}

func handlerReset(s *state.State, _ Command) error {
	return s.Db.ResetUsers(context.Background())
}

func handlerUsers(s *state.State, _ Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		line := fmt.Sprintf("* %s", user.Name)
		if user.Name == s.Config.CurrentUser {
			line += " (current)"
		}
		fmt.Printf("%s\n", line)
	}
	return nil
}

func handlerAggregate(_ *state.State, _ Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command needs name of feed!")
	}
	if len(cmd.Args) == 1 {
		return fmt.Errorf("command needs URL to the feed!")
	}
	feedname := cmd.Args[0]
	url := cmd.Args[1]
	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedname,
		Url:       url,
		UserID:    user.ID})
	if err != nil {
		return err
	}
	fmt.Println("feed created:", feed.Name)
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID})
	if err != nil {
		return err
	}
	fmt.Println("feed", feedFollow.FeedName, "is now followed by", feedFollow.UserName)
	return nil
}

func handlerFeeds(s *state.State, _ Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("name: %s\turl: %s\tcreated by: %s\n", feed.Name, feed.Url, feed.Creator)
	}
	return nil
}

func handlerFollow(s *state.State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command needs url!")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID})
	if err != nil {
		return err
	}
	fmt.Println(feedFollow.UserName, "is now following", feedFollow.FeedName)
	return nil
}

func handlerFollows(s *state.State, _ Command, currentUser database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}
	fmt.Println(currentUser.Name, "is following:")
	for _, f := range feedFollows {
		fmt.Println("\t", f.FeedName)
	}
	return nil
}

func withUserLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}
