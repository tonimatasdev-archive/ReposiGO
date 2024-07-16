package repo

import "os"

const (
	Public  = "PUBLIC"
	Secret  = "SECRET"
	Private = "PRIVATE"
)

type Repository struct {
	Name    string
	Id      string
	Type    string
	Primary bool
}

func (repo Repository) GetName() string {
	return repo.Name
}

func RepositoryInit(name string, id string, repoType string, primary bool) Repository {
	repo := Repository{name, id, repoType, primary}

	_ = os.MkdirAll("repositories/"+id, 0755)

	return repo
}
