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
    Title       string          `json:"title" bson:",omitempty"`
    Summary     string          `json:"summary" bson:",omitempty"`
    Publisher   user            `json:"user,omitempty" bson:",omitempty"`
    File        bson.ObjectId   `json:"file,omitempty" bson:",omitempty"`
}

var MediaCollection = configs.MongoDB.C("media")
var MediaGridFS = configs.MongoDB.GridFS("media")