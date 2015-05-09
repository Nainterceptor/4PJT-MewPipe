package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
)

type user struct {
    *User
}

type Media struct {
    Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    Title       string          `json:"title"`
    Summary     string          `json:"summary"`
    Publisher   user            `json:"user,omitempty" bson:"user,omitempty"`
    File        bson.ObjectId   `json:"file,omitempty" bson:"file,omitempty"`
}

var MediaCollection = configs.MongoDB.C("media")
var MediaGridFS = configs.MongoDB.GridFS("media")