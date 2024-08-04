package session

import (
	"encoding/base64"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"net"
	"net/http"
	"strings"
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
		if parts[0] == value.Username && parts[1] == value.Token {
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
