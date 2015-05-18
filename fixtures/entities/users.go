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

func getFooUser() *entities.User {
	usr := entities.UserNew()
	usr.Email = "foo@bar.com"
	usr.Password = "foobar"
	usr.Name.FirstName = "Foo"
	usr.Name.LastName = "Bar"
	usr.Name.NickName = "FooBar"
	return usr
}

func getAdminUser() *entities.User {
	usr := entities.UserNew()
	usr.Email = "admin@admin.com"
	usr.Password = "admin"
	usr.Name.FirstName = "Admin"
	usr.Name.LastName = "Admin"
	usr.Name.NickName = "Admin"
	usr.Roles = append(usr.Roles, "Admin")
	return usr
}

func InsertSomeUser() {
	var userArray []*entities.User
	userArray = append(userArray, getFooUser())
	userArray = append(userArray, getAdminUser())
	for _, usr := range userArray {
		if err := usr.Insert(); err != nil {
			panic(err)
			return
		}
	}

	fmt.Println("Users added")
}
