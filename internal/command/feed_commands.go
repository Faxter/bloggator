package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/state"
	"github.com/google/uuid"
)

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
