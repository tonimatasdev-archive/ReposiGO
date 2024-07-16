package session

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

type Session struct {
	Username    string
	Token       string
	ReadAccess  []string
	WriteAccess []string
}

func CreateSession(username string, readAccess []string, writeAccess []string) {
	token, err := generateRandomToken(64)

	if err != nil {
		log.Println("Error creating the session.")
	} else {
		value := sessions[username]

		if value.Username == username {
			log.Println("Session \"" + username + "\" already exists.")
			return
		}

		sessions[username] = Session{username, token, readAccess, writeAccess}
		log.Println("Session created successfully.")
	}
}

func DeleteSession(username string) {
	value := sessions[username]

	if value.Username == "" {
		log.Println("Session", "\""+username+"\"", "not found.")
	} else {
		delete(sessions, username)
	}
}

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
