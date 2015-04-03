package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
)

type User struct {
    Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name     string        `json:"name"`
}

var UserCollection = configs.MongoDB.C("users")