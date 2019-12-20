package pullrequest

// Registry .
type Registry interface {
	VCSRepository() GitRepository
}

// Configuration .
type Configuration interface {
	GithubToken() string
	BaseBranch() string
}
