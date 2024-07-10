package main

import (
	"container/list"
	"encoding/base64"
	"fmt"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var repositories = list.New()

const (
	username = "test"
	password = "dev-version"
)

func main() {
	releaseRepository := repo.RepositoryInit("Releases", "releases", repo.Public, true)
	secretRepository := repo.RepositoryInit("Secret", "secret", repo.Secret, false)
	privateRepository := repo.RepositoryInit("Private", "private", repo.Private, false)

	repositories.PushFront(releaseRepository)
	repositories.PushFront(secretRepository)
	repositories.PushFront(privateRepository)

	for e := repositories.Front(); e != nil; e = e.Next() {
		value, ok := e.Value.(repo.Repository)

		if !ok {
			continue
		}

		http.HandleFunc("/"+value.Id+"/", auth(value))

		if value.Primary {
			http.HandleFunc("/", auth(value))
		}
	}

	fmt.Println("Server listening on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func auth(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if repository.Type == repo.Private || r.Method == http.MethodPut {
			if !checkAuth(r.Header.Get("Authorization")) {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		handleRequest(w, r, repository)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, repository repo.Repository) {
	switch r.Method {
	case http.MethodPut:
		handlePut(w, r, repository)
	case http.MethodGet:
		handleGet(w, r, repository)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func checkAuth(auth string) bool {
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

	return parts[0] == username && parts[1] == password
}

func handlePut(w http.ResponseWriter, r *http.Request, repository repo.Repository) {
	filePath := getFilePath(r, repository)

	if filePath == "" {
		http.NotFound(w, r)
	}

	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleGet(w http.ResponseWriter, r *http.Request, repository repo.Repository) {
	filePath := getFilePath(r, repository)

	if filePath == "" {
		http.NotFound(w, r)
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	http.ServeFile(w, r, filePath)
}

func getFilePath(r *http.Request, repository repo.Repository) string {
	if strings.Contains(r.URL.Path, "..") {
		return ""
	}

	if repository.Primary {
		return "." + repository.Id + r.URL.Path
	} else {
		return "." + r.URL.Path
	}
}
