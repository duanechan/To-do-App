package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(key)
}
