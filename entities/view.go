package entities

import (
	"supinfo/mewpipe/configs"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var viewCollection = configs.MongoDB.C("media.views")

type View struct {
	Id    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User  bson.ObjectId `json:"user,omitempty" bson:"user,omitempty"`
	Media bson.ObjectId `json:"media,omitempty" bson:"media,omitempty"`
	Count int           `json:"count" bson:"count"`
}

var change = mgo.Change{
	Update:    bson.M{"$inc": bson.M{"count": 1}},
	Upsert:    true,
	ReturnNew: true,
}

func ViewNewAnonymous(media bson.ObjectId) error {
	return upsertOnCriteria(bson.M{"media": media})
}

func ViewNew(user bson.ObjectId, media bson.ObjectId) error {
	return upsertOnCriteria(bson.M{"user": user, "media": media})
}

func upsertOnCriteria(query interface{}) error {
	view := new(View)
	_, err := viewCollection.Find(query).Apply(change, &view)
	return err
}
