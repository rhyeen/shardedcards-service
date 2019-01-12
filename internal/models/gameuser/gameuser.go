package gameuser

import "gopkg.in/mgo.v2/bson"

// User are the details of a gameuser.
type User struct {
	ObjectID bson.ObjectId `json:"-" bson:"_id,omitempty"`
	ID       string        `json:"id" bson:"id"`
}
