package session

import (
	"encoding/base64"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"net/http"
	"strings"
)

func CheckAuth(sessions []Session, auth string, r *http.Request, repository repo.Repository) bool {
	if auth == "" {
		return false
	}

	authParts := strings.SplitN(auth, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return false
	}

	for _, value := range sessions {
		if parts[0] == value.Username && parts[1] == value.Token {
			if r.Method == http.MethodPut && (utils.Contains(value.WriteAccess, repository.Id) || utils.Contains(value.ReadAccess, "*")) {
				return true
			} else if repository.Type == repo.Private && (utils.Contains(value.ReadAccess, repository.Id) || utils.Contains(value.ReadAccess, "*")) {
				return true
			}
		}
	}

	return false
}
