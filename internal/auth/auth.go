package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	var t *jwt.Token // create variable t of type *jwt.Token

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{ //create header and payload
		Audience:  jwt.ClaimStrings{username},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Second)),
	})

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing key")
		return signedToken, err
	}

	return signedToken, nil
}
