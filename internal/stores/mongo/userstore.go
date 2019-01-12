package mongo

import (
	"github.com/pkg/errors"
	"github.com/rhyeen/shardedcards-service/internal/models/gameuser"
	"github.com/rhyeen/shardedcards-service/internal/stores/storeerror"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const usersCollection = "users"

// UserStore is the mongo store for users
type UserStore struct {
	mongoDB *mgo.Database
}

// NewUserStore returns a UserStore
func NewUserStore(mongoDB *mgo.Database) UserStore {
	return UserStore{
		mongoDB: mongoDB,
	}
}

// EnsureIndices ensures the mongo indices on the collection
func (s UserStore) EnsureIndices() error {
	indices := []mgo.Index{
		// @TODO
	}
	err := EnsureIndicesOnCollection(s.mongoDB, usersCollection, indices)
	if err != nil {
		return errors.Wrapf(err, "failed to update %v collection", usersCollection)
	}
	return nil
}

// CreateUser creates a new user.
func (s UserStore) CreateUser(item gameuser.User) (gameuser.User, error) {
	item.ObjectID = bson.NewObjectId()
	db := s.mongoDB.Session.Copy().DB(s.mongoDB.Name)
	defer db.Session.Close()
	err := db.C(usersCollection).Insert(&item)
	if mgo.IsDup(err) {
		return item, &storeerror.DupEntry{
			ID:  item.ObjectID.String(), // @TODO: should be the GUI
			Err: err,
		}
	}
	return item, err
}
