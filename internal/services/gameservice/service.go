package gameservice

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rhyeen/shardedcards-service/internal/models/game"
	"github.com/rhyeen/shardedcards-service/internal/models/gameuser"
	"github.com/rhyeen/shardedcards-service/internal/stores/store"
)

// GameService is the service for handling game-related APIs
type GameService struct {
	GameStore store.GameStore
	DeckStore store.DeckStore
	UserStore store.UserStore
}

// CreateGameParams params for CreateGame
type CreateGameParams struct {
	User gameuser.User
	Type string
}

// CreateGame creates a new game.
func (s GameService) CreateGame(ctx context.Context, params CreateGameParams) (game.Game, error) {
	game, err := s.GameStore.CreateGame(game.Game{})
	if err != nil {
		return game, errors.Wrapf(err, "failed to create game: %+v", params)
	}
	return game, nil
}
