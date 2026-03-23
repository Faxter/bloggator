package command

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/faxter/bloggator/internal/database"
	"github.com/faxter/bloggator/internal/rss"
	"github.com/faxter/bloggator/internal/state"
	"github.com/google/uuid"
)

func handlerAggregate(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command needs time between requests!")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if scrapeFeeds(s) != nil {
			fmt.Println("scraping failed", err)
			return err
		}
	}
}

func scrapeFeeds(s *state.State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feed, err = s.Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	fmt.Println("found RSS channel:", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
		pubDate, timeParseErr := time.Parse(time.RFC3339, item.PubDate)
		_, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: timeParseErr == nil},
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			fmt.Printf("error creating post %s: %s\n", item.Link, err)
		}
	}
	return nil
}

func handlerBrowse(s *state.State, cmd Command, currentUser database.User) error {
	postLimit := 2
	if len(cmd.Args) > 0 {
		var err error
		postLimit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("could not convert input to integer: %w", err)
		}
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit:  int32(postLimit)})
	if err != nil {
		return fmt.Errorf("could not retrieve posts from database: %w", err)
	}
	for _, post := range posts {
		fmt.Println(post.Title)
		fmt.Println("-------------------")
		fmt.Println(post.Description.String)
		fmt.Printf("\n\n\n")
	}
	return nil
}
