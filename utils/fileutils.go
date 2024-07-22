package utils

import (
	"github.com/TonimatasDEV/ReposiGO/repo"
	"net/http"
	"strings"
)

func FilePath(r *http.Request, repository repo.Repository, primary repo.Repository) string {
	if strings.Contains(r.URL.Path, "../") {
		return ""
	}

	if repository == primary {
		return "repositories/" + repository.Id + r.URL.Path
	} else {
		return "repositories/" + r.URL.Path
	}
}
