package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/sawadashota/orb-update/internal/configfile"
)

// UpdateAll orbs
func (h *Handler) UpdateAll() error {
	cfs := make([]*configfile.ConfigFile, 0, len(h.c.TargetFiles()))
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

	updates, err := configfile.DetectUpdateSet(cfs)
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		return nil
	}

	for _, update := range updates {
		if err := h.Update(cfs, update); err != nil {
			return err
		}
	}

	return nil
}

// Update an orb
func (h *Handler) Update(cfs []*configfile.ConfigFile, diff *configfile.Update) error {
	ctx := context.Background()

	if h.doesCreatePullRequest {
		alreadyCreated, err := h.r.PullRequest().AlreadyCreated(ctx, h.branchForPR(diff))
		if err != nil {
			return err
		}

		if alreadyCreated {
			_, _ = fmt.Fprintf(h.r.Logger(), "PR for %s has been already created\n", diff.After.String())
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
		diff.After.Namespace(),
		diff.After.Name(),
		diff.Before.Version(),
		diff.After.Version(),
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

func (h *Handler) branchForPR(diff *configfile.Update) string {
	return fmt.Sprintf("%s/%s/%s-%s", h.c.GitBranchPrefix(), diff.After.Namespace(), diff.After.Name(), diff.After.Version())
}

func commitMessage(diff *configfile.Update) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Before.Namespace(), diff.Before.Name(), diff.Before.Version(), diff.After.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Before.Namespace(), diff.Before.Name())
	return message
}
