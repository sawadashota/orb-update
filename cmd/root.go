package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sawadashota/orb-update/configfile"

	"github.com/pkg/errors"
	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/filesystem"
	"github.com/sawadashota/orb-update/git"
	"github.com/sawadashota/orb-update/orb"
	"github.com/sawadashota/orb-update/pullrequest"
	"github.com/spf13/cobra"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			d := driver.NewDefaultDriver()

			g, fs, err := newGitRepository(d, repo)
			if err != nil {
				return err
			}

			reader, err := fs.Reader(filePath)
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
				if err := update(d, g, fs, cf, diff, repo, filePath, doesCreatePullRequest); err != nil {
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

func update(d driver.Driver, g git.Git, fs filesystem.Filesystem, cf *configfile.ConfigFile, diff *orb.Difference, repo, filePath string, doesCreatePullRequest bool) error {
	var pr pullrequest.Creator
	ctx := context.Background()

	if doesCreatePullRequest {
		r, err := parseRepository(repo)
		if err != nil {
			return err
		}

		pr, err = pullrequest.NewGitHubPullRequest(ctx, d, r.owner, r.name, diff)
		if err != nil {
			return err
		}

		alreadyCreated, err := pr.AlreadyCreated(ctx, branchForPR(diff))
		if err != nil {
			return err
		}

		if alreadyCreated {
			_, _ = fmt.Fprintf(os.Stdout, "PR for %s has been already created\n", diff.New.String())
			return nil
		}

		if err := g.Switch(branchForPR(diff), true); err != nil {
			return err
		}
		defer func() {
			if err := g.SwitchBack(); err != nil {
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
	writer, err := fs.OverWriter(filePath)
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

	if _, err := g.Commit(commitMessage(diff), filePath); err != nil {
		return err
	}

	if err := g.Push(ctx, branchForPR(diff)); err != nil {
		return err
	}

	if err := pr.Create(ctx, commitMessage(diff), branchForPR(diff)); err != nil {
		return err
	}

	return nil
}
func newGitRepository(d driver.Driver, repo string) (git.Git, filesystem.Filesystem, error) {
	if d.Configuration().FilesystemStrategy() == configuration.OsFileSystemStrategy {
		g, fs, err := git.OpenCurrentDirectoryRepository(d)
		if err != nil {
			return nil, nil, err
		}

		return g, fs, nil
	}

	if repo == "" {
		return nil, nil, errors.New("repository name wasn't given")
	}

	r, err := parseRepository(repo)
	if err != nil {
		return nil, nil, err
	}

	g, fs, err := git.Clone(d, r.owner, r.name)
	if err != nil {
		return nil, nil, err
	}

	return g, fs, nil
}

func branchForPR(diff *orb.Difference) string {
	return fmt.Sprintf("orb-update/%s/%s-%s", diff.New.Namespace(), diff.New.Name(), diff.New.Version())
}

func commitMessage(diff *orb.Difference) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Old.Namespace(), diff.Old.Name(), diff.Old.Version(), diff.New.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Old.Namespace(), diff.Old.Name())
	return message
}
