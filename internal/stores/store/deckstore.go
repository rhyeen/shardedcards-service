package store

import "github.com/rhyeen/shardedcards-service/internal/models/deck"

// DeckStore defines the required functionality for any associated store.
type DeckStore interface {
	CreateDeck(item deck.Deck) (deck.Deck, error)
}
