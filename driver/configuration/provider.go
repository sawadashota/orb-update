package configuration

// Provider of configuration
type Provider interface {
	GitAuthorName() string
	GitAuthorEmail() string

	GitHubPullRequest() bool
	GithubUsername() string
	GithubToken() string

	BaseBranch() string

	FilesystemStrategy() string

	RepositoryName() string
	FilePath() string
}
