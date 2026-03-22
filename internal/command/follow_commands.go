package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/state"
	"github.com/google/uuid"
)

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

func handlerUnfollow(s *state.State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command needs url!")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.Db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(currentUser.Name, "unfollowed", feed.Name)
	return nil
}
