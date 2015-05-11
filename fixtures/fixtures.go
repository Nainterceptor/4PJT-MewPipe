package main

import (
    "supinfo/mewpipe/fixtures/entities"
    "fmt"
)

func main() {
    Users()
}

func Users() {
    fmt.Println("Users Fixtures")
    entities.ClearUsers()
    entities.InsertSomeUser()
}