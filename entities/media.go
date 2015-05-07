package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
)


type Media struct {
    Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    Title       string          `json:"title"`
    Summary     string          `json:"summary"`
    Publisher   User   `json:"user"`
    File        bson.ObjectId   `json:"file,omitempty"`
}

var MediaCollection = configs.MongoDB.C("media")
var MediaGridFS = configs.MongoDB.GridFS("media")