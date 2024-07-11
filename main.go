package main

import (
	"bufio"
	"container/list"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/session"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var repositories = list.New()
var sessions = list.New()
var primaryRepository repo.Repository

func main() {
	releaseRepository := repo.RepositoryInit("Releases", "releases", repo.Public, true)
	secretRepository := repo.RepositoryInit("Secret", "secret", repo.Secret, false)
	privateRepository := repo.RepositoryInit("Private", "private", repo.Private, false)

	primaryRepository = releaseRepository
	repositories.PushFront(secretRepository)
	repositories.PushFront(privateRepository)

	http.HandleFunc("/", handleRequest)

	session, createUserErr := session.SessionInit("test", []string{"*"}, []string{"*"})
	if createUserErr != nil {
		log.Println("Error creating the session:", createUserErr)
	} else {
		log.Println(session.Username)
		log.Println(session.Token)
		sessions.PushFront(session)
	}

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

	log.Println("Server listening on port 8080")

	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-stopChan:
			stop(server)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		rawCommand, err := inputReader.ReadString('\n')

		command := strings.Replace(rawCommand, "\n", "", -1)

		if err != nil {
			log.Println("Exception on read the command:", err)
		}

		if command == "exit" || command == "stop" {
			stop(server)
		}

		log.Println("Command:", command)
	}
}

func stop(server *http.Server) {
	log.Println("ReposiGO stopped successfully.")
	_ = server.Close()
	os.Exit(0)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	found := false
	var repository repo.Repository

	for e := repositories.Front(); e != nil; e = e.Next() {
		value, ok := e.Value.(repo.Repository)

		if !ok {
			continue
		}

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
		if !session.CheckAuth(sessions, r.Header.Get("Authorization"), r, repository) {
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
	defer file.Close()

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
	defer file.Close()

	http.ServeFile(w, r, filePath)
}
