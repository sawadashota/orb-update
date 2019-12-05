package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sawadashota/orb-update/orb"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var filePath string

	c := &cobra.Command{
		Use:     "orb-update",
		Short:   "Update CircleCI Orb versions",
		Example: "",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}

			path := filepath.Join(pwd, filePath)
			file, err := os.Stat(path)
			if err != nil {
				return err
			}

			if file.IsDir() {
				return errors.Errorf("%s is not a file", filePath)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}

			path := filepath.Join(pwd, filePath)

			reader, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				return err
			}

			cf, err := orb.NewConfigFile(reader)
			if err != nil {
				return err
			}
			reader.Close()

			conf, err := cf.Parse()
			if err != nil {
				return err
			}

			cl := orb.NewDefaultClient()
			newVersions := make([]*orb.Orb, 0, len(conf.Orbs))
			for _, o := range conf.Orbs {
				newVersion, err := cl.LatestVersion(o)
				if err != nil {
					return err
				}

				if o.Version() != newVersion.Version() {
					_, _ = fmt.Fprintf(os.Stdout, "Updating %s/%s (%s => %s)\n", o.Namespace(), o.Name(), o.Version(), newVersion.Version())
				}

				newVersions = append(newVersions, newVersion)
			}

			// overwrite the config file
			writer, err := os.Create(path)
			if err != nil {
				return err
			}
			defer writer.Close()

			if len(newVersions) == 0 {
				return nil
			}

			return cf.Update(writer, newVersions...)
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".circleci/config.yml", "target config file path")

	return c
}
