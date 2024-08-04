package main

import (
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"github.com/TonimatasDEV/ReposiGO/console"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/session"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	config, err := configuration.LoadConfig()

	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	configuration.ServerConfig = *config

	for _, configRepository := range config.Repositories {
		repository := repo.RepositoryInit(configRepository.Name, configRepository.Id, configRepository.Type)

		if repository.Id == config.Primary {
			repo.PrimaryRepository = repository
		} else {
			repo.Repositories = append(repo.Repositories, repository)
		}
	}

	if repo.PrimaryRepository.Id != config.Primary {
		log.Fatal("Primary repository not found.")
	}

	http.HandleFunc("/", handleRequest)

	portStr := strconv.Itoa(config.Port)
	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: nil,
	}

	go func() {
		var err error

		if config.CertFile == "none" || config.KeyFile == "none" {
			err = server.ListenAndServe()
		} else {
			err = server.ListenAndServeTLS(config.CertFile, config.KeyFile)
		}

		if err != nil {
			log.Fatal("Server starting error: ", err)
			return
		}
	}()

	log.Println("Server listening on port " + portStr + ".")

	go session.BanHandler()

	console.Console(server)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	found := false
	var repository repo.Repository

	for _, value := range repo.Repositories {
		if strings.HasPrefix(r.URL.Path, "/"+value.Id+"/") {
			found = true
			repository = value
			break
		}
	}

	if !found {
		repository = repo.PrimaryRepository
	}

	if repository.Type == repo.Private || r.Method == http.MethodPut {
		result, str, code := session.CheckAuth(r.Header.Get("Authorization"), r, repository)

		if !result {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, str, code)
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
