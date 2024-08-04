package repo

import (
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
