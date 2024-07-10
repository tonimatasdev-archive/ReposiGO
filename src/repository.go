package src

type Repository struct {
	Name string
	Id   string
}

func (repo Repository) getName() string {
	return repo.Name
}

func RepositoryInit(name string, id string) Repository {
	repo := Repository{name, id}
	return repo
}
