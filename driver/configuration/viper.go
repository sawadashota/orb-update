package configuration

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	viperGitAuthorName      = "git.author.name"
	viperGitAuthorEmail     = "git.author.email"
	viperGithubToken        = "github.token"
	viperGithubUsername     = "github.username"
	viperBaseBranch         = "base_branch"
	viperFilesystemStrategy = "filesystem.strategy"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// ViperProvider .
type ViperProvider struct{}

// NewViperProvider .
func NewViperProvider() *ViperProvider {
	return new(ViperProvider)
}

// GitAuthorName .
func (v *ViperProvider) GitAuthorName() string {
	return viper.GetString(viperGitAuthorName)
}

// GitAuthorEmail .
func (v *ViperProvider) GitAuthorEmail() string {
	return viper.GetString(viperGitAuthorEmail)
}

// GithubUsername.
func (v *ViperProvider) GithubUsername() string {
	return viper.GetString(viperGithubUsername)
}

// GithubToken .
func (v *ViperProvider) GithubToken() string {
	return viper.GetString(viperGithubToken)
}

// BaseBranch .
func (v *ViperProvider) BaseBranch() string {
	branch := viper.GetString(viperBaseBranch)
	if branch == "" {
		return "master"
	}
	return branch
}

const (
	InMemoryFilesystemStrategy = "memory"
	OsFileSystemStrategy       = "os"
)

// FilesystemStrategy .
func (v *ViperProvider) FilesystemStrategy() string {
	strategy := viper.GetString(viperFilesystemStrategy)
	if strategy == "" {
		return OsFileSystemStrategy
	}

	return InMemoryFilesystemStrategy
}
