package configuration

type Provider interface {
	GitAuthorName() string
	GitAuthorEmail() string

	GithubUsername() string
	GithubToken() string

	TargetBranch() string

	FilesystemStrategy() string
}
