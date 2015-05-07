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
    Name        name            `json:"name"`
    Email       string          `json:"email"`
	HashedPassword    string    `json:"hashedpassword"`
    UserTokens  []UserToken     `json:"usertokens"`
}

type UserToken struct {
    Token   string `json:"token"`
}

type Registration struct {
	*User
	Password    string    `json:"password"`
}

type Connexion struct {
    Email       string    `json:"email"`
    Password    string    `json:"password"`
}

var UserCollection = configs.MongoDB.C("users")
