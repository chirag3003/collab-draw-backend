package repository

var repo *Repository

type Repository struct {
	Project   ProjectRepository
	Workspace WorkspaceRepository
}

func Setup() *Repository {
	repo = &Repository{
		Project:   NewProjectRepository(),
		Workspace: NewWorkspaceRepository(),
	}
	return repo
}
