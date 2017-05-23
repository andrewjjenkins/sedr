package main

import (
	"client"
	"errors"
	"fmt"
	"os"
)

func getLoginCredentials() (email string, password string, err error) {
	emailEnv, emailPresent := os.LookupEnv("SEDR_EMAIL")
	passwordEnv, passwordPresent := os.LookupEnv("SEDR_PASSWORD")

	if !emailPresent || !passwordPresent {
		return "", "", errors.New("No credentials. Specify Elite: Dangerous email " +
			"and password via environment variables SEDR_EMAIL and SEDR_PASSWORD")
	}
	return emailEnv, passwordEnv, nil
}

func main() {
	fmt.Println("Hello world")
	c, err := client.NewEDClient()
	if err != nil {
		panic(err)
	}

	email, password, err := getLoginCredentials()
	if err != nil {
		panic(err)
	}
	err = c.Login(email, password)
	if err != nil {
		panic(err)
	}
}
