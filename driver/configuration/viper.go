package configuration

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	viperGithubToken    = "github.token"
	viperGithubUsername = "github.username"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

type ViperProvider struct{}

func NewViperProvider() *ViperProvider {
	return new(ViperProvider)
}

func (v *ViperProvider) GithubUsername() string {
	return viper.GetString(viperGithubUsername)
}

func (v *ViperProvider) GithubToken() string {
	return viper.GetString(viperGithubToken)
}
