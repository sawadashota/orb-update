package extraction

import (
	"github.com/sawadashota/orb-update/internal/orb"
)

//type UpdateDetection struct {
//	r io.Reader
//}
//
//func (ud *UpdateDetection) Do() ([]*Update, error) {
//
//}

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
func (cf *Extraction) DetectUpdate() ([]*Update, error) {
	orbs, err := cf.Do()
	if err != nil {
		return nil, err
	}

	cl := orb.NewDefaultClient()
	updates := newUpdateSet()
	for _, o := range orbs {
		if !o.Version().IsSemantic() {
			continue
		}

		newVersion, err := cl.LatestVersion(o)
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
