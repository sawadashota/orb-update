package configuration

type Provider interface {
	GithubUsername() string
	GithubToken() string
}
