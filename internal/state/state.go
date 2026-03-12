package state

import (
	"github.com/faxter/bloggator/internal/config"
	"github.com/faxter/bloggator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

func NewState(db *database.Queries, cfg *config.Config) *State {
	return &State{Db: db, Config: cfg}
}
