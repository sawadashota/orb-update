package vcsuser

// Configuration .
type Configuration interface {
	GithubUsername() string
	GithubToken() string
	BaseBranch() string
}
