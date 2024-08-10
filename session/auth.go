package session

import (
	"encoding/base64"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"strings"
	"time"
)

func CheckAuth(auth string, r *http.Request, repository repo.Repository) (bool, string, int) {
	if auth == "" {
		return false, "Incomplete auth", http.StatusUnauthorized
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return false, "Internal Server Error (Ip)", http.StatusUnauthorized
	}

	isBanned, str := checkBan(ip)

	if isBanned {
		return false, str, http.StatusTooManyRequests
	}

	authParts := strings.SplitN(auth, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return false, "Incomplete auth", http.StatusUnauthorized
	}

	decoded, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		return false, "Internal Server Error (Decoding)", http.StatusUnauthorized
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return false, "Incomplete auth", http.StatusBadRequest
	}

	for _, value := range sessions {
		passwordErr := comparePasswordWithDelay(value.HashedToken, parts[1])

		if parts[0] == value.Username && passwordErr == nil {
			if r.Method == http.MethodPut && (utils.Contains(value.WriteAccess, repository.Id) || utils.Contains(value.ReadAccess, "*")) {
				return true, "", 0
			} else if repository.Type == repo.Private && (utils.Contains(value.ReadAccess, repository.Id) || utils.Contains(value.ReadAccess, "*")) {
				return true, "", 0
			}
		}
	}

	addTry(ip)

	return false, "Unauthorized", http.StatusUnauthorized
}

func comparePasswordWithDelay(token string, input string) error {
	startTime := time.Now()

	passwordErr := bcrypt.CompareHashAndPassword([]byte(token), []byte(input))

	delay := 100*time.Millisecond - time.Since(startTime)

	if delay > 0 {
		time.Sleep(delay)
	}

	return passwordErr
}
