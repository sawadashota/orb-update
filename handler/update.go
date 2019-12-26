package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sawadashota/orb-update/internal/extraction"
)

// UpdateAll orbs
func (h *Handler) UpdateAll() error {
	files := make([]io.Reader, 0, len(h.c.TargetFiles()))

	var err error
	for i, target := range h.c.TargetFiles() {
		files[i], err = h.r.Filesystem().Reader(target)
		if err != nil {
			return err
		}
	}

	reader := io.MultiReader(files...)
	cf, err := extraction.New(reader)
	if err != nil {
		return err
	}

	updates, err := cf.DetectUpdate()
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		return nil
	}

	for _, update := range updates {
		if err := h.Update(cf, update); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) filterOrbs(updates []*extraction.Update) []*extraction.Update {
	var filtered []*extraction.Update
	for _, update := range updates {
		for _, ignoreOrb := range h.c.TargetFiles() {
			if strings.HasPrefix(update.Before.String(), ignoreOrb) {
				continue
			}
		}
		filtered = append(filtered, update)
	}
	return filtered
}

// Update an orb
func (h *Handler) Update(cf *extraction.Extraction, diff *extraction.Update) error {
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

	for _, filePath := range h.c.TargetFiles() {
		writer, err := h.r.Filesystem().OverWriter(filePath)
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

func (h *Handler) branchForPR(diff *extraction.Update) string {
	return fmt.Sprintf("%s/%s/%s-%s", h.c.GitBranchPrefix(), diff.After.Namespace(), diff.After.Name(), diff.After.Version())
}

func commitMessage(diff *extraction.Update) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Before.Namespace(), diff.Before.Name(), diff.Before.Version(), diff.After.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Before.Namespace(), diff.Before.Name())
	return message
}
