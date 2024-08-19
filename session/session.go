package session

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var sessions = make(map[string]Session)

type Session struct {
	Username    string   `json:"username"`
	HashedToken string   `json:"hashedToken"`
	ReadAccess  []string `json:"readAccess"`
	WriteAccess []string `json:"writeAccess"`
}

func CreateSession(username string, readAccess []string, writeAccess []string) {
	token, err := generateRandomToken(54)

	hashedToken, err1 := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if err != nil || err1 != nil {
		log.Println("Error creating the session", err, err1)
	} else {
		value := sessions[username]

		if value.Username == username {
			log.Println("Session \"" + username + "\" already exists.")
			return
		}

		sessions[username] = Session{username, string(hashedToken), readAccess, writeAccess}
		saveSessions()

		log.Println("Session \"" + username + "\" created successfully with the token \"" + token + "\".")
	}
}

func DeleteSession(username string) {
	value := sessions[username]

	if value.Username == "" {
		log.Println("Session", "\""+username+"\"", "not found.")
	} else {
		delete(sessions, username)
		saveSessions()
		log.Println("Session", "\""+username+"\"", "deleted.")
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
