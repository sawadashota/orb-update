package configuration

import (
	"github.com/spf13/viper"
)

const (
	// ViperGitAuthorName can be set from outside when empty
	ViperGitAuthorName = "git.author.name"
	// ViperGitAuthorEmail can be set from outside when empty
	ViperGitAuthorEmail = "git.author.email"

	viperGitBranchPrefix   = "git.branch_prefix"
	defaultGitBranchPrefix = "orb-update"

	viperRepositoryName     = "repository.name"
	viperTargetFiles        = "target_files"
	viperPullRequest        = "github.pull_request"
	viperGithubToken        = "github.token"
	viperGithubUsername     = "github.username"
	viperBaseBranch         = "base_branch"
	viperFilesystemStrategy = "filesystem.strategy"
)

// ViperProvider .
type ViperProvider struct{}

// NewViperProvider .
func NewViperProvider() *ViperProvider {
	return new(ViperProvider)
}

// GitBranchPrefix .
func (v *ViperProvider) GitBranchPrefix() string {
	prefix := viper.GetString(viperGitBranchPrefix)
	if prefix == "" {
		return defaultGitBranchPrefix
	}

	return prefix
}

// GitAuthorName .
func (v *ViperProvider) GitAuthorName() string {
	return viper.GetString(ViperGitAuthorName)
}

// GitAuthorEmail .
func (v *ViperProvider) GitAuthorEmail() string {
	return viper.GetString(ViperGitAuthorEmail)
}

// GitHubPullRequest .
func (v *ViperProvider) GitHubPullRequest() bool {
	return viper.GetBool(viperPullRequest)
}

// GithubUsername .
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
	// InMemoryFilesystemStrategy .
	InMemoryFilesystemStrategy = "memory"
	// OsFileSystemStrategy .
	OsFileSystemStrategy = "os"
)

// FilesystemStrategy .
func (v *ViperProvider) FilesystemStrategy() string {
	strategy := viper.GetString(viperFilesystemStrategy)
	if strategy == "" {
		return OsFileSystemStrategy
	}

	return strategy
}

// RepositoryName .
func (v *ViperProvider) RepositoryName() string {
	return viper.GetString(viperRepositoryName)
}

// TargetFiles .
func (v *ViperProvider) TargetFiles() []string {
	return viper.GetStringSlice(viperTargetFiles)
}
