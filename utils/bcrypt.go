package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword[:])
}
