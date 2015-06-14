package entities

import (
	"supinfo/mewpipe/configs"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getShareCountCollection() *mgo.Collection {
	return configs.MongoDB.C("media.shareCounts")
}

type ShareCount struct {
	Id    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User  bson.ObjectId `json:"user,omitempty" bson:"user,omitempty"`
	Media bson.ObjectId `json:"media,omitempty" bson:"media,omitempty"`
	Count int           `json:"count" bson:"count"`
}

var changeShareCount = mgo.Change{
	Update:    bson.M{"$inc": bson.M{"count": 1}},
	Upsert:    true,
	ReturnNew: true,
}

func ShareCountNewAnonymous(media bson.ObjectId) error {
	return upsertShareCountOnCriteria(bson.M{"media": media})
}

func ShareCountNew(user bson.ObjectId, media bson.ObjectId) error {
	return upsertShareCountOnCriteria(bson.M{"user": user, "media": media})
}

func upsertShareCountOnCriteria(query interface{}) error {
	view := new(View)
	_, err := getShareCountCollection().Find(query).Apply(changeShareCount, &view)
	return err
}
