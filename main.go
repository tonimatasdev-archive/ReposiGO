package main

import (
	"github.com/TonimatasDEV/ReposiGO/console"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/session"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	repositories      []repo.Repository
	primaryRepository repo.Repository
)

func main() {
	releaseRepository := repo.RepositoryInit("Releases", "releases", repo.Public, true)
	secretRepository := repo.RepositoryInit("Secret", "secret", repo.Secret, false)
	privateRepository := repo.RepositoryInit("Private", "private", repo.Private, false)

	primaryRepository = releaseRepository
	repositories = append(repositories, secretRepository)
	repositories = append(repositories, privateRepository)

	http.HandleFunc("/", handleRequest)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	log.Println("Server listening on port 8080.")

	console.Console(server)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	found := false
	var repository repo.Repository

	for _, value := range repositories {
		if strings.HasPrefix(r.URL.Path, "/"+value.Id+"/") {
			found = true
			repository = value
			break
		}
	}

	if !found {
		repository = primaryRepository
	}

	if repository.Type == repo.Private || r.Method == http.MethodPut {
		if !session.CheckAuth(r.Header.Get("Authorization"), r, repository) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	switch r.Method {
	case http.MethodPut:
		handlePut(w, r, repository)
	case http.MethodGet:
		handleGet(w, r, repository)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePut(w http.ResponseWriter, r *http.Request, repository repo.Repository) {
	filePath := utils.FilePath(r, repository)

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
	defer utils.FileError(file)

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleGet(w http.ResponseWriter, r *http.Request, repository repo.Repository) {
	filePath := utils.FilePath(r, repository)

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
	defer utils.FileError(file)

	http.ServeFile(w, r, filePath)
}
