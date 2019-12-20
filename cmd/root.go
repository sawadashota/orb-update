package cmd

import (
	"github.com/sawadashota/orb-update/driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd .
func RootCmd() *cobra.Command {
	var config string

	c := &cobra.Command{
		Use:     "orb-update",
		Short:   "Update CircleCI Orb versions",
		Example: "",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AddConfigPath(config)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			d, err := driver.NewDefaultDriver()
			if err != nil {
				return err
			}

			return d.Registry().Handler().UpdateAll()
		},
	}

	c.Flags().StringVarP(&config, "config", "c", ".orb-update.yml", "configuration file for orb-update")

	return c
}
