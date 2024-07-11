package utils

import (
	"fmt"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"net/http"
	"strings"
)

func FilePath(r *http.Request, repository repo.Repository) string {
	if strings.Contains(r.URL.Path, "../") {
		return ""
	}

	if repository.Primary {
		fmt.Println("." + repository.Id + r.URL.Path)
		return "./" + repository.Id + r.URL.Path
	} else {
		return "." + r.URL.Path
	}
}
