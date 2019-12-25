package configfile

import (
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

// hasUpdate or not
func (u *Update) hasUpdate() bool {
	return u.Before.Version() != u.After.Version()
}

// DetectUpdateSet reads multi config file and differences set
func DetectUpdateSet(cfs []*ConfigFile) ([]*Update, error) {
	set := NewUpdateSet()
	for _, cf := range cfs {
		diffs, err := cf.DetectUpdate()
		if err != nil {
			return nil, err
		}

		set.addMulti(diffs...)
	}

	return set.set, nil
}

// updateSet is set of Update
// set should not be multiply
type updateSet struct {
	set []*Update
}

// NewUpdateSet .
func NewUpdateSet() *updateSet {
	return &updateSet{set: make([]*Update, 0)}
}

// addMulti updates
func (ds *updateSet) addMulti(updates ...*Update) {
	for _, diff := range updates {
		ds.add(diff)
	}
}

// add update
func (ds *updateSet) add(update *Update) {
	for _, d := range ds.set {
		if d.Before.IsSameOrb(update.Before) {
			return
		}
	}
	ds.set = append(ds.set, update)
}

// DetectUpdate from CircleCI config file
func (cf *ConfigFile) DetectUpdate() ([]*Update, error) {
	orbs, err := cf.ExtractOrbs()
	if err != nil {
		return nil, err
	}

	cl := orb.NewDefaultClient()
	updates := make([]*Update, 0, len(orbs))
	for _, o := range orbs {
		if !o.Version().IsSemantic() {
			continue
		}

		newVersion, err := cl.LatestVersion(o)
		if err != nil {
			return nil, err
		}

		diff := NewUpdate(o, newVersion)
		if !diff.hasUpdate() {
			continue
		}

		if o.Version() != newVersion.Version() {
			updates = append(updates, diff)
		}
	}

	return updates, nil
}
