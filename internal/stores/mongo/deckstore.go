package mongo

import (
	"github.com/pkg/errors"
	"github.com/rhyeen/shardedcards-service/internal/models/deck"
	"github.com/rhyeen/shardedcards-service/internal/stores/storeerror"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const decksCollection = "decks"

// DeckStore is the mongo store for decks
type DeckStore struct {
	mongoDB *mgo.Database
}

// NewDeckStore returns a DeckStore
func NewDeckStore(mongoDB *mgo.Database) DeckStore {
	return DeckStore{
		mongoDB: mongoDB,
	}
}

// EnsureIndices ensures the mongo indices on the collection
func (s DeckStore) EnsureIndices() error {
	indices := []mgo.Index{
		// @TODO
	}
	err := EnsureIndicesOnCollection(s.mongoDB, decksCollection, indices)
	if err != nil {
		return errors.Wrapf(err, "failed to update %v collection", decksCollection)
	}
	return nil
}

// CreateDeck creates a new deck.
func (s DeckStore) CreateDeck(item deck.Deck) (deck.Deck, error) {
	item.ObjectID = bson.NewObjectId()
	db := s.mongoDB.Session.Copy().DB(s.mongoDB.Name)
	defer db.Session.Close()
	err := db.C(decksCollection).Insert(&item)
	if mgo.IsDup(err) {
		return item, &storeerror.DupEntry{
			ID:  item.ObjectID.String(), // @TODO: should be the GUI
			Err: err,
		}
	}
	return item, err
}
