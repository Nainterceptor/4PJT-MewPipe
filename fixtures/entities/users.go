package entities

import (
    "supinfo/mewpipe/entities"
    "fmt"
    "gopkg.in/mgo.v2/bson"
    "supinfo/mewpipe/utils"
)

const (
    FirstName = "John"
    LastName = "Beau"
    Email = "JBeau@gmail.com"
    Password = "test"
)

func ClearUsers() {

    usr := entities.User{}
    if _, err := entities.UserCollection.RemoveAll(&usr); err != nil {
        panic(err)
        return
    }
    fmt.Println("All users deleted")
}

func InsertBasicUser() {

    usr := entities.User{}
    usr.Id = bson.NewObjectId()
    usr.Name.FirstName = FirstName
    usr.Name.LastName = LastName
    usr.Email = Email
    usr.Password = Password
    usr.HashedPassword = utils.Hash(usr.Password)
    usr.Password = ""

    if err := entities.UserCollection.Insert(&usr); err != nil {
        panic(err)
        return
    }

    fmt.Println("User add: ", Email)
}
