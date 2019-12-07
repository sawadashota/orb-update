package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sawadashota/orb-update/driver"

	"github.com/sawadashota/orb-update/pullrequest"

	"github.com/sawadashota/orb-update/orb"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var filePath string
	var repo string

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

			diffs, err := parse(cf, filePath)
			if err != nil {
				return err
			}

			if len(diffs) == 0 {
				return nil
			}

			// overwrite the configuration file
			writer, err := os.Create(path)
			if err != nil {
				return err
			}
			defer writer.Close()

			for _, diff := range diffs {
				_, _ = fmt.Fprintf(
					os.Stdout,
					"Updating %s/%s (%s => %s)\n",
					diff.New.Namespace(),
					diff.New.Name(),
					diff.Old.Version(),
					diff.New.Version(),
				)

				if err := cf.Update(writer, diff.New); err != nil {
					return err
				}

				if repo == "" {
					continue
				}

				d := driver.NewDefaultDriver()

				r, err := parseRepository(repo)
				if err != nil {
					return err
				}

				ctx := context.Background()
				pr, err := pullrequest.NewGitHubPullRequest(ctx, d, r.owner, r.name, diff)
				if err != nil {
					return err
				}

				message := fmt.Sprintf("Bump %s/%s from %s to %s\n", diff.Old.Namespace(), diff.Old.Name(), diff.Old.Version(), diff.New.Version())
				message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Old.Namespace(), diff.Old.Name())
				if err := pr.Create(ctx, message); err != nil {
					return err
				}
			}

			return nil
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".circleci/config.yml", "target config file path")
	c.Flags().StringVarP(&repo, "repo", "r", "", "GitHub Repository to create Pull Request")

	return c
}
