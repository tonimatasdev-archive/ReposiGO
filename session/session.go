package session

import (
	"crypto/rand"
	"encoding/base64"
)

type Session struct {
	Username    string
	Token       string
	ReadAccess  []string
	WriteAccess []string
}

func SessionInit(username string, readAccess []string, writeAccess []string) (Session, error) {
	token, err := generateRandomToken(64)
	session := Session{username, token, readAccess, writeAccess}

	if err != nil {
		return session, err
	}

	return session, nil
}

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
