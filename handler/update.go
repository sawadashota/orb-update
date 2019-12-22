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
	cfs := make([]orb.ConfigFile, 0, len(h.c.TargetFiles()))
	for _, target := range h.c.TargetFiles() {
		reader, err := h.r.Filesystem().Reader(target)
		if err != nil {
			return err
		}

		cf, err := configfile.New(reader, target)
		if err != nil {
			return err
		}
		reader.Close()

		cfs = append(cfs, cf)
	}

	diffs, err := orb.DetectUpdateSet(cfs)
	if err != nil {
		return err
	}

	if len(diffs) == 0 {
		return nil
	}

	for _, diff := range diffs {
		if err := h.Update(cfs, diff); err != nil {
			return err
		}
	}

	return nil
}

// Update an orb
func (h *Handler) Update(cfs []orb.ConfigFile, diff *orb.Difference) error {
	ctx := context.Background()

	if h.doesCreatePullRequest {
		alreadyCreated, err := h.r.PullRequest().AlreadyCreated(ctx, h.branchForPR(diff))
		if err != nil {
			return err
		}

		if alreadyCreated {
			_, _ = fmt.Fprintf(h.r.Logger(), "PR for %s has been already created\n", diff.New.String())
			return nil
		}

		if err := h.r.Git().Switch(h.branchForPR(diff), true); err != nil {
			return err
		}
		defer func() {
			if err := h.r.Git().SwitchBack(); err != nil {
				_, _ = fmt.Fprintln(os.Stdout, err)
			}
		}()
	}

	_, _ = fmt.Fprintf(
		h.r.Logger(),
		"Updating %s/%s (%s => %s)\n",
		diff.New.Namespace(),
		diff.New.Name(),
		diff.Old.Version(),
		diff.New.Version(),
	)

	for _, cf := range cfs {
		// overwrite the configuration file
		writer, err := h.r.Filesystem().OverWriter(cf.Path())
		if err != nil {
			return err
		}

		if err := cf.Update(writer, diff); err != nil {
			return err
		}
		writer.Close()
	}

	if !h.doesCreatePullRequest {
		return nil
	}

	if _, err := h.r.Git().Commit(commitMessage(diff), h.c.TargetFiles()); err != nil {
		return err
	}

	if err := h.r.Git().Push(ctx, h.branchForPR(diff)); err != nil {
		return err
	}

	if err := h.r.PullRequest().Create(ctx, diff, commitMessage(diff), h.branchForPR(diff)); err != nil {
		return err
	}

	return nil
}

func (h *Handler) branchForPR(diff *orb.Difference) string {
	return fmt.Sprintf("%s/%s/%s-%s", h.c.GitBranchPrefix(), diff.New.Namespace(), diff.New.Name(), diff.New.Version())
}

func commitMessage(diff *orb.Difference) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Old.Namespace(), diff.Old.Name(), diff.Old.Version(), diff.New.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Old.Namespace(), diff.Old.Name())
	return message
}
