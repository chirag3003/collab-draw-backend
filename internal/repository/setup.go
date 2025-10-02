package repository

var repo *Repository

type Repository struct {
	Project ProjectRepository
}

func Setup() *Repository {
	repo = &Repository{
		Project: NewProjectRepository(),
	}
	return repo
}
