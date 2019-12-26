package configuration

// Provider of configuration
type Provider interface {
	GitBranchPrefix() string
	GitAuthorName() string
	GitAuthorEmail() string

	GitHubPullRequest() bool
	GithubUsername() string
	GithubToken() string

	BaseBranch() string

	FilesystemStrategy() string

	RepositoryName() string
	TargetFiles() []string
	IgnoreOrbs() []string
}
