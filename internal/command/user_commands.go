package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/state"
	"github.com/google/uuid"
)

func withUserLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}

func handlerReset(s *state.State, _ Command) error {
	return s.Db.ResetUsers(context.Background())
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
