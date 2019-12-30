package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sawadashota/orb-update/driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd of cobra command
var RootCmd = &cobra.Command{
	Use:   "orb-update",
	Short: "Update CircleCI Orb versions",
	RunE: func(_ *cobra.Command, _ []string) error {
		d, err := driver.NewDefaultDriver()
		if err != nil {
			return err
		}

		return d.Registry().Handler().UpdateAll()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// configuration file path for orb-update
var config string

func initConfig() {
	viper.SetConfigFile(config)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `config was file not found because "%s"`, err)
		_, _ = fmt.Fprintln(os.Stderr, "")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", ".orb-update.yml", "configuration file path for orb-update")
}
