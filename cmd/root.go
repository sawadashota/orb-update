package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sawadashota/orb-update/configfile"
	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/orb"
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

			reader, err := d.Registry().Filesystem().Reader(filePath)
			if err != nil {
				return err
			}

			cf, err := configfile.New(reader)
			if err != nil {
				return err
			}
			reader.Close()

			diffs, err := differences(cf)
			if err != nil {
				return err
			}

			if len(diffs) == 0 {
				return nil
			}

			for _, diff := range diffs {
				if err := update(d, cf, diff, doesCreatePullRequest); err != nil {
					return err
				}
			}

			return nil
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".circleci/config.yml", "target config file path")
	c.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository name ex) owner/name")
	c.Flags().BoolVarP(&doesCreatePullRequest, "pull-request", "p", false, "Create Pull Request or not")

	return c
}

func update(d driver.Driver, cf *configfile.ConfigFile, diff *orb.Difference, doesCreatePullRequest bool) error {
	ctx := context.Background()

	if doesCreatePullRequest {
		alreadyCreated, err := d.Registry().PullRequest().AlreadyCreated(ctx, branchForPR(diff))
		if err != nil {
			return err
		}

		if alreadyCreated {
			_, _ = fmt.Fprintf(os.Stdout, "PR for %s has been already created\n", diff.New.String())
			return nil
		}

		if err := d.Registry().Git().Switch(branchForPR(diff), true); err != nil {
			return err
		}
		defer func() {
			if err := d.Registry().Git().SwitchBack(); err != nil {
				_, _ = fmt.Fprintln(os.Stdout, err)
			}
		}()
	}

	_, _ = fmt.Fprintf(
		os.Stdout,
		"Updating %s/%s (%s => %s)\n",
		diff.New.Namespace(),
		diff.New.Name(),
		diff.Old.Version(),
		diff.New.Version(),
	)

	// overwrite the configuration file
	writer, err := d.Registry().Filesystem().OverWriter(d.Configuration().FilePath())
	if err != nil {
		return err
	}

	if err := cf.Update(writer, diff); err != nil {
		return err
	}
	writer.Close()

	if !doesCreatePullRequest {
		return nil
	}

	if _, err := d.Registry().Git().Commit(commitMessage(diff), d.Configuration().FilePath()); err != nil {
		return err
	}

	if err := d.Registry().Git().Push(ctx, branchForPR(diff)); err != nil {
		return err
	}

	if err := d.Registry().PullRequest().Create(ctx, diff, commitMessage(diff), branchForPR(diff)); err != nil {
		return err
	}

	return nil
}

func branchForPR(diff *orb.Difference) string {
	return fmt.Sprintf("orb-update/%s/%s-%s", diff.New.Namespace(), diff.New.Name(), diff.New.Version())
}

func commitMessage(diff *orb.Difference) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Old.Namespace(), diff.Old.Name(), diff.Old.Version(), diff.New.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Old.Namespace(), diff.Old.Name())
	return message
}
