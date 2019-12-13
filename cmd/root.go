package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sawadashota/orb-update/driver"
	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/filesystem"
	"github.com/sawadashota/orb-update/git"
	"github.com/sawadashota/orb-update/orb"
	"github.com/sawadashota/orb-update/pullrequest"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var filePath string
	var repo string
	var branch string
	var doesCreatePullRequest bool

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
			d := driver.NewDefaultDriver()

			var fs filesystem.Filesystem
			var g git.Git
			var err error
			if d.Configuration().FilesystemStrategy() == configuration.InMemoryFilesystemStrategy {
				if repo == "" {
					return errors.New("repository name wasn't given")
				}

				r, err := parseRepository(repo)
				if err != nil {
					return err
				}

				g, fs, err = git.Clone(d, r.owner, r.name, branch)
				if err != nil {
					return err
				}
			} else {
				g, fs, err = git.OpenCurrentDirectoryRepository(d)
				if err != nil {
					return err
				}
			}

			//pwd, err := os.Getwd()
			//if err != nil {
			//	return err
			//}
			//
			//path := filepath.Join(pwd, filePath)

			reader, err := fs.Reader(filePath)
			//reader, err := os.OpenFile(path, os.O_RDONLY, 0666)
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

			for _, diff := range diffs {
				err = func() error {
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
						defer g.SwitchBack()
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

					if err := cf.Update(writer, diff.New); err != nil {
						return err
					}
					writer.Close()

					if !doesCreatePullRequest {
						return nil
					}

					if err := pr.Create(ctx, commitMessage(diff), branchForPR(diff)); err != nil {
						return err
					}

					if _, err := g.Commit(commitMessage(diff), branchForPR(diff)); err != nil {
						return err
					}

					if err := g.Push(ctx, branch); err != nil {
						return err
					}

					return nil
				}()

				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".circleci/config.yml", "target config file path")
	c.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository name ex) owner/name")
	c.Flags().StringVarP(&branch, "branch", "b", "master", "Branch to clone")
	c.Flags().BoolVarP(&doesCreatePullRequest, "pull-request", "p", false, "Create Pull Request or not")

	return c
}

func branchForPR(diff *orb.Difference) string {
	return fmt.Sprintf("orb-update/%s/%s-%s", diff.New.Namespace(), diff.New.Name(), diff.New.Version())
}

func commitMessage(diff *orb.Difference) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n", diff.Old.Namespace(), diff.Old.Name(), diff.Old.Version(), diff.New.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Old.Namespace(), diff.Old.Name())
	return message
}
