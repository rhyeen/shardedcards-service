package store

import "github.com/rhyeen/shardedcards-service/internal/models/gameuser"

// UserStore defines the required functionality for any associated store.
type UserStore interface {
	CreateUser(item gameuser.User) (gameuser.User, error)
}
