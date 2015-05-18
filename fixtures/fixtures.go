package main

import (
	"fmt"
	"os"
	"supinfo/mewpipe/configs"
	"supinfo/mewpipe/fixtures/entities"
)

func main() {
	configs.Parse()
	var args []string = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Full import")
		fmt.Println("===========")
		Users()
		return
	}
	if contains(args, "users") {
		fmt.Println("User import")
		fmt.Println("===========")
		Users()
		return
	}
	fmt.Println("Undefined import")
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func Users() {
	entities.ClearUsers()
	entities.InsertSomeUser()
}
