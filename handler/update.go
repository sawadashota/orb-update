package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/sawadashota/orb-update/internal/configfile"
	"github.com/sawadashota/orb-update/internal/orb"
)

// UpdateAll orbs
func (h *Handler) UpdateAll() error {
	reader, err := h.r.Filesystem().Reader(h.c.FilePath())
	if err != nil {
		return err
	}

	cf, err := configfile.New(reader)
	if err != nil {
		return err
	}
	reader.Close()

	diffs, err := orb.DetectUpdate(cf)
	if err != nil {
		return err
	}

	if len(diffs) == 0 {
		return nil
	}

	for _, diff := range diffs {
		if err := h.Update(cf, diff); err != nil {
			return err
		}
	}

	return nil
}

// Update an orb
func (h *Handler) Update(cf *configfile.ConfigFile, diff *orb.Difference) error {
	ctx := context.Background()

	if h.doesCreatePullRequest {
		alreadyCreated, err := h.r.PullRequest().AlreadyCreated(ctx, branchForPR(diff))
		if err != nil {
			return err
		}

		if alreadyCreated {
			_, _ = fmt.Fprintf(h.logger, "PR for %s has been already created\n", diff.New.String())
			return nil
		}

		if err := h.r.Git().Switch(branchForPR(diff), true); err != nil {
			return err
		}
		defer func() {
			if err := h.r.Git().SwitchBack(); err != nil {
				_, _ = fmt.Fprintln(os.Stdout, err)
			}
		}()
	}

	_, _ = fmt.Fprintf(
		h.logger,
		"Updating %s/%s (%s => %s)\n",
		diff.New.Namespace(),
		diff.New.Name(),
		diff.Old.Version(),
		diff.New.Version(),
	)

	// overwrite the configuration file
	writer, err := h.r.Filesystem().OverWriter(h.c.FilePath())
	if err != nil {
		return err
	}

	if err := cf.Update(writer, diff); err != nil {
		return err
	}
	writer.Close()

	if !h.doesCreatePullRequest {
		return nil
	}

	if _, err := h.r.Git().Commit(commitMessage(diff), h.c.FilePath()); err != nil {
		return err
	}

	if err := h.r.Git().Push(ctx, branchForPR(diff)); err != nil {
		return err
	}

	if err := h.r.PullRequest().Create(ctx, diff, commitMessage(diff), branchForPR(diff)); err != nil {
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
