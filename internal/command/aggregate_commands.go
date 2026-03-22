package command

import (
	"context"
	"fmt"
	"time"

	"github.com/faxter/bloggator/internal/rss"
	"github.com/faxter/bloggator/internal/state"
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
	}
	return nil
}
