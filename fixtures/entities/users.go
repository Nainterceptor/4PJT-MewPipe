package entities

import (
	"fmt"
	"supinfo/mewpipe/configs"
	"supinfo/mewpipe/entities"
)

func ClearUsers() {
	if err := configs.MongoDB.C("users").DropCollection(); err != nil {
		panic(err)
		return
	}
	fmt.Println("All users deleted")
}

func InsertSomeUser() {
	usr := entities.UserNew()
	usr.Email = "foo@bar.com"
	usr.Password = "foobar"
	usr.Name.FirstName = "Foo"
	usr.Name.LastName = "Bar"
	usr.Name.NickName = "FooBar"

	if err := usr.Insert(); err != nil {
		panic(err)
		return
	}

	fmt.Println("Users added")
}
