package deck

import (
	"github.com/rhyeen/shardedcards-service/internal/models/gameuser"
	"gopkg.in/mgo.v2/bson"
)

// Deck is a user's deck of cards.
type Deck struct {
	ObjectID bson.ObjectId `json:"-" bson:"_id,omitempty"`
	User     gameuser.User `json:"user" bson:"user"`
}
