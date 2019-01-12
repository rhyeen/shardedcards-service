package mongo

import (
	"github.com/pkg/errors"
	"github.com/rhyeen/shardedcards-service/internal/models/game"
	"github.com/rhyeen/shardedcards-service/internal/stores/storeerror"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const gamesCollection = "games"

// GameStore is the mongo store for games
type GameStore struct {
	mongoDB *mgo.Database
}

// NewGameStore returns a GameStore
func NewGameStore(mongoDB *mgo.Database) GameStore {
	return GameStore{
		mongoDB: mongoDB,
	}
}

// EnsureIndices ensures the mongo indices on the collection
func (s GameStore) EnsureIndices() error {
	indices := []mgo.Index{
		// @TODO
	}
	err := EnsureIndicesOnCollection(s.mongoDB, gamesCollection, indices)
	if err != nil {
		return errors.Wrapf(err, "failed to update %v collection", gamesCollection)
	}
	return nil
}

// CreateGame creates a new game.
func (s GameStore) CreateGame(item game.Game) (game.Game, error) {
	item.ObjectID = bson.NewObjectId()
	db := s.mongoDB.Session.Copy().DB(s.mongoDB.Name)
	defer db.Session.Close()
	err := db.C(gamesCollection).Insert(&item)
	if mgo.IsDup(err) {
		return item, &storeerror.DupEntry{
			ID:  item.ObjectID.String(), // @TODO: should be the GUI
			Err: err,
		}
	}
	return item, err
}
