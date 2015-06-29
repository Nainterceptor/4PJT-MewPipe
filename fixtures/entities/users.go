package entities

import (
	"fmt"

	"github.com/Nainterceptor/4PJT-MewPipe/configs"
	"github.com/Nainterceptor/4PJT-MewPipe/entities"

	"gopkg.in/mgo.v2/bson"
)

func ClearUsers() {
	configs.MongoDB.C("users").DropCollection()
	fmt.Println("All users deleted")
	fmt.Println("Admin user is admin@admin.com with 'Admin' password.")
}

func getFooUser() *entities.User {
	usr := entities.UserNewFromId(bson.ObjectIdHex("555a076a2fd06c1891000001"))
	usr.Email = "foo@bar.com"
	usr.Password = "foobar"
	usr.Name.FirstName = "Foo"
	usr.Name.LastName = "Bar"
	usr.Name.NickName = "FooBar"
	return usr
}

func getAdminUser() *entities.User {
	usr := entities.UserNewFromId(bson.ObjectIdHex("555a076a2fd06c1891000002"))
	usr.Email = "admin@admin.com"
	usr.Password = "admin"
	usr.Name.FirstName = "Admin"
	usr.Name.LastName = "Admin"
	usr.Name.NickName = "Admin"
	usr.Roles = append(usr.Roles, "Admin")
	return usr
}

func getTwitterUser() *entities.User {
	usr := entities.UserNew()
	usr.Name.FirstName = "Twitter"
	usr.Name.LastName = "Twitter"
	usr.Name.NickName = "Twitter"
	usr.Twitter.UserId = "999999"
	return usr
}

func InsertSomeUser() {
	var userArray []*entities.User
	userArray = append(userArray, getFooUser())
	userArray = append(userArray, getAdminUser())
	userArray = append(userArray, getTwitterUser())
	for _, usr := range userArray {
		if err := usr.Insert(); err != nil {
			panic(err)
			return
		}
		usr.TokenNew()
	}

	fmt.Println("Users added")
}
