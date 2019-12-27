package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sawadashota/orb-update/internal/extraction"
)

// UpdateAll orbs
func (h *Handler) UpdateAll() error {
	e, err := h.r.Extraction()
	if err != nil {
		return err
	}

	orbFilters := extraction.ExcludeMatchPackages(h.c.IgnoreOrbs())
	updates, err := e.Updates(orbFilters...)
	if err != nil {
		return err
	}

	for _, update := range updates {
		if err := h.Update(e, update); err != nil {
			return err
		}
	}

	return nil
}

// Update an orb
func (h *Handler) Update(e *extraction.Extraction, update *extraction.Update) error {
	ctx := context.Background()

	alreadyCreated, switchBack, err := h.beforeUpdate(ctx, update)
	if err != nil {
		return err
	}
	if alreadyCreated {
		_, _ = fmt.Fprintf(h.r.Logger(), "PR for %s has been already created\n", update.After)
		return nil
	}
	defer switchBack()

	_, _ = fmt.Fprintf(
		h.r.Logger(),
		"Updating %s/%s (%s => %s)\n",
		update.After.Namespace(),
		update.After.Name(),
		update.Before.Version(),
		update.After.Version(),
	)

	for _, filePath := range h.c.TargetFiles() {
		writer, err := h.r.Filesystem().OverWriter(filePath)
		if err != nil {
			return err
		}

		var b bytes.Buffer

		scan := bufio.NewScanner(e.Reader())
		for scan.Scan() {
			func() {
				if strings.Contains(scan.Text(), update.Before.String()) {
					b.WriteString(
						extraction.OrbFormatRegex.ReplaceAllString(scan.Text(),
							"$1@"+update.After.Version().String()),
					)
					b.WriteString("\n")
					return
				}
				b.Write(scan.Bytes())
				b.WriteString("\n")
			}()
		}

		_, err = io.Copy(writer, &b)
		if err != nil {
			return err
		}
		writer.Close()
	}

	if !h.doesCreatePullRequest {
		return nil
	}

	if _, err := h.r.Git().Commit(commitMessage(update), h.c.TargetFiles()); err != nil {
		return err
	}

	if err := h.r.Git().Push(ctx, h.branchForPR(update)); err != nil {
		return err
	}

	if err := h.r.PullRequest().Create(ctx, update, commitMessage(update), h.branchForPR(update)); err != nil {
		return err
	}

	return nil
}

func (h *Handler) beforeUpdate(ctx context.Context, update *extraction.Update) (alreadyCreated bool, switchBack func(), err error) {
	if !h.doesCreatePullRequest {
		return
	}

	alreadyCreated, err = h.r.PullRequest().AlreadyCreated(ctx, h.branchForPR(update))
	if err != nil {
		return
	}

	if alreadyCreated {
		_, _ = fmt.Fprintf(h.r.Logger(), "PR for %s has been already created\n", update.After)
		return
	}

	if err = h.r.Git().Switch(h.branchForPR(update), true); err != nil {
		return
	}
	switchBack = func() {
		if err := h.r.Git().SwitchBack(); err != nil {
			_, _ = fmt.Fprintln(os.Stdout, err)
		}
	}
	return
}

func (h *Handler) branchForPR(diff *extraction.Update) string {
	return fmt.Sprintf("%s/%s/%s-%s", h.c.GitBranchPrefix(), diff.After.Namespace(), diff.After.Name(), diff.After.Version())
}

func commitMessage(diff *extraction.Update) string {
	message := fmt.Sprintf("Bump %s/%s from %s to %s\n\n", diff.Before.Namespace(), diff.Before.Name(), diff.Before.Version(), diff.After.Version())
	message += fmt.Sprintf("https://circleci.com/orbs/registry/orb/%s/%s", diff.Before.Namespace(), diff.Before.Name())
	return message
}
