package cmd

import (
	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd .
func RootCmd() *cobra.Command {
	var filePath string
	var repo string
	var doesCreatePullRequest bool

	c := &cobra.Command{
		Use:     "orb-update",
		Short:   "Update CircleCI Orb versions",
		Example: "",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if repo != "" {
				viper.Set(configuration.ViperRepositoryName, repo)
			}

			viper.Set(configuration.ViperFilePath, filePath)

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			d, err := driver.NewDefaultDriver()
			if err != nil {
				return err
			}

			return d.Registry().Handler().UpdateAll()
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".circleci/config.yml", "target config file path")
	c.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository name ex) owner/name")
	c.Flags().BoolVarP(&doesCreatePullRequest, "pull-request", "p", false, "Create Pull Request or not")

	return c
}
