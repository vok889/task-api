package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func main() {
	password, err := hashPassword("secret")
	if err != nil {
		fmt.Println("error=", err.Error)
	}
	fmt.Println("password=", password)
}
