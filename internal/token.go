package internal

import (
	"errors"
	"log"
	"os"
)

func GetToken() (string, error) {
	tkn := os.Getenv("GITHUB_TOKEN")
	if tkn == "" {
		log.Print("Token not found. You must set it in your environment like")
		log.Print("export GITHUB_TOKEN=000a0aaaa0000a00000000aaa00000000a000000")
		log.Print("You can generate a token at https://github.com/settings/tokens")
		return "", errors.New("No token available in system")
	}

	return tkn, nil
}
