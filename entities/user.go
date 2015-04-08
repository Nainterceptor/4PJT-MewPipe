package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
)

type name struct {
    Firstname   string  `json:"firstname"`
    Lastname   string  `json:"lastname"`
}

type User struct {
    Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    Name        name            `json:"name"`
    Email       string          `json:"email"`
	HashedPassword    string    `json:"hashedpassword"`
}

type Registration struct {
	*User
	Password    string    `json:"password"`
}

var UserCollection = configs.MongoDB.C("users")
