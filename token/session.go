package token

import (
	"crypto/rand"
	"encoding/base64"
)

type Session struct {
	Username string
	Token    string
}

func SessionInit(username string) (Session, error) {
	token, err := generateRandomToken(64)
	session := Session{username, token}

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
