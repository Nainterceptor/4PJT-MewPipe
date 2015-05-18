package main

import (
	"fmt"
	"supinfo/mewpipe/fixtures/entities"
)

func main() {
	Users()
}

func Users() {
	fmt.Println("Users Fixtures")
	entities.ClearUsers()
	entities.InsertSomeUser()
}
