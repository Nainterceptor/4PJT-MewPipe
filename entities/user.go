package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
)

type name struct {
    FirstName   string  `json:"firstname"`
    LastName   string  `json:"lastname"`
}

type User struct {
    Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    Name        name            `json:"name" bson:",omitempty"`
    Email       string          `json:"email" bson:",omitempty"`
    Password    string    `json:"password" bson:",omitempty"`
    HashedPassword    string    `json:"hashedpassword" bson:",omitempty"`
    UserTokens  []UserToken     `json:"usertokens" bson:",omitempty"`
}

type UserToken struct {
    Token   string `json:"token"`
}

var UserCollection = configs.MongoDB.C("users")
