package main

import (
	"errors"
	"fmt"
	"github.com/andrewjjenkins/sedr/client"
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
	c, err := client.NewEDClient()
	if err != nil {
		panic(err)
	}

	if c.NeedLogin() {
		email, password, err := getLoginCredentials()
		if err != nil {
			panic(err)
		}
		err = c.Login(email, password)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Using stored cookie\n")
	}
	if err != nil {
		panic(err)
	}
}
