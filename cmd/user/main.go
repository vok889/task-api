package main

import (
	"fmt"
	"task-api/internal/user"
)

func main() {
	password, err := user.HashPassword("secret")
	if err != nil {
		fmt.Println("error=", err.Error)
	}
	fmt.Println("password=", password)
}
