package store

import "github.com/rhyeen/shardedcards-service/internal/models/game"

// GameStore defines the required functionality for any associated store.
type GameStore interface {
	CreateGame(item game.Game) (game.Game, error)
}
