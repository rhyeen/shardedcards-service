package game

import (
	"github.com/rhyeen/shardedcards-service/internal/models/gameuser"
	"gopkg.in/mgo.v2/bson"
)

// Type is a valid game type.
type Type string

// All the valid values for Type
const (
	TypeDungeonDelve Type = "dungeonDelve"
)

// Game is a user's game session
type Game struct {
	ObjectID bson.ObjectId `json:"-" bson:"_id,omitempty"`
	ID       string        `json:"id" bson:"id"`
	User     gameuser.User `json:"user" bson:"user"`
	Type     Type          `json:"type" bson:"type"`
}
