package main

import (
    "supinfo/mewpipe/fixtures/entities"
    "os"
    "fmt"
)

func main() {
    if (os.Args[1] == "user") {
        Users()
        fmt.Println("Users Fixtures")
    } else {
        fmt.Println("'", os.Args[1], "'", "is not support")
    }
}

func Users() {
    entities.ClearUsers()
    entities.InsertBasicUser()
}