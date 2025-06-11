package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		panic("JWT_SECRET must be at least 32 characters")
	}

	rolesFlag := flag.String("roles", "", "Comma-separated list of roles (required)")
	flag.Parse()

	if *rolesFlag == "" {
		fmt.Println("Error: --roles flag is required")
		flag.Usage()
		os.Exit(1)
	}

	roles := strings.Split(*rolesFlag, ",")
	if len(roles) == 0 || (len(roles) == 1 && roles[0] == "") {
		fmt.Println("Error: at least one role must be provided")
		os.Exit(1)
	}

	claims := CustomClaims{
		UserID: 1,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	fmt.Println("Bearer", tokenString)
}
