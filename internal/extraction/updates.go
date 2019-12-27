package extraction

import (
	"strings"

	"github.com/sawadashota/orb-update/internal/orb"
)

// Update of version between new and old
type Update struct {
	Before *orb.Orb
	After  *orb.Orb
}

// NewUpdate .
func NewUpdate(before, after *orb.Orb) *Update {
	return &Update{
		Before: before,
		After:  after,
	}
}

// validate or not
func (u *Update) validate() bool {
	return u.Before.Version() != u.After.Version()
}

// updateSet is set of Update
// set should not be multiply
type updateSet struct {
	set []*Update
}

// newUpdateSet .
func newUpdateSet() *updateSet {
	return &updateSet{set: make([]*Update, 0)}
}

// add update
func (us *updateSet) add(update *Update) {
	for _, d := range us.set {
		if d.Before.IsSameOrb(update.Before) {
			return
		}
	}
	us.set = append(us.set, update)
}

// Match with orb function
type Match func(o *orb.Orb) bool

// MatchPackage matcher
func MatchPackage(name string) Match {
	return func(o *orb.Orb) bool {
		return strings.HasPrefix(o.String(), name)
	}
}

// Filter orbs function
type Filter func(orbs []*orb.Orb) []*orb.Orb

// Exclude matched orbs
func Exclude(match Match) Filter {
	return func(orbs []*orb.Orb) []*orb.Orb {
		filtered := make([]*orb.Orb, 0, len(orbs))
		for _, o := range orbs {
			if !match(o) {
				filtered = append(filtered, o)
			}
		}
		return filtered
	}
}

// ExcludeMatchPackages combines Exclude and MatchPackage functions
func ExcludeMatchPackages(names []string) []Filter {
	filters := make([]Filter, 0, len(names))
	for _, name := range names {
		filter := Exclude(MatchPackage(name))
		filters = append(filters, filter)
	}
	return filters
}

// Updates from CircleCI config file
func (e *Extraction) Updates(filters ...Filter) ([]*Update, error) {
	orbs, err := e.Orbs()
	if err != nil {
		return nil, err
	}

	for _, filter := range filters {
		orbs = filter(orbs)
	}

	updates := newUpdateSet()
	for _, o := range orbs {
		if !o.Version().IsSemantic() {
			continue
		}

		newVersion, err := DefaultClient.LatestVersion(o)
		if err != nil {
			return nil, err
		}

		update := NewUpdate(o, newVersion)
		if !update.validate() {
			continue
		}

		if o.Version() != newVersion.Version() {
			updates.add(update)
		}
	}

	return updates.set, nil
}
