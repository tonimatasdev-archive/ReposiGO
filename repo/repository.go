package repo

import (
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"log"
	"os"
)

const (
	Public  = "PUBLIC"
	Secret  = "SECRET"
	Private = "PRIVATE"
)

var (
	Repositories      []Repository
	PrimaryRepository Repository
)

type Repository struct {
	Name string
	Id   string
	Type string
}

func (repo Repository) GetName() string {
	return repo.Name
}

func RepositoryInit(name string, id string, repoType string) Repository {
	repo := Repository{name, id, repoType}

	_ = os.MkdirAll("repositories/"+id, 0755)

	return repo
}

func InitRepositories() {
	primary := configuration.ServerConfig.Primary

	for _, configRepository := range configuration.ServerConfig.Repositories {
		repository := RepositoryInit(configRepository.Name, configRepository.Id, configRepository.Type)

		if repository.Id == primary {
			PrimaryRepository = repository
		} else {
			Repositories = append(Repositories, repository)
		}
	}

	if PrimaryRepository.Id != primary {
		log.Fatal("Primary repository not found.")
	}
}
