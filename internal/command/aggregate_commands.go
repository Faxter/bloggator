package command

import (
	"context"
	"fmt"

	"github.com/faxter/bloggator/internal/rss"
	"github.com/faxter/bloggator/internal/state"
)

func handlerAggregate(_ *state.State, _ Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
