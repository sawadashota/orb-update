package orb

import (
	"io"
)

// Difference of version between new and old
type Difference struct {
	Old *Orb
	New *Orb
}

// NewDifference .
func NewDifference(old, new *Orb) *Difference {
	return &Difference{
		Old: old,
		New: new,
	}
}

// HasUpdate or not
func (d *Difference) HasUpdate() bool {
	return d.Old.Version() != d.New.Version()
}

// ConfigFile of CircleCI
type ConfigFile interface {
	Parse() ([]*Orb, error)
	Path() string
	Update(w io.Writer, diff *Difference) error
}

// DetectUpdateFromMultipleFile reads multi config file and differences set
func DetectUpdateFromMultipleFile(cfs []ConfigFile) ([]*Difference, error) {
	multiDiffs := make([]*Difference, 0)
	for _, cf := range cfs {
		diffs, err := DetectUpdate(cf)
		if err != nil {
			return nil, err
		}

		multiDiffs = append(multiDiffs, diffs...)
	}

	multiDiffSet := make([]*Difference, 0)
	for _, diff := range multiDiffs {
		if !hasOrb(diff, multiDiffSet) {
			multiDiffSet = append(multiDiffSet, diff)
		}
	}

	return multiDiffSet, nil
}

func hasOrb(needle *Difference, haystack []*Difference) bool {
	for _, d := range haystack {
		if needle.Old.Namespace() == d.Old.Namespace() && needle.Old.Name() == d.Old.Name() {
			return true
		}
	}

	return false
}

// DetectUpdate from CircleCI config file
func DetectUpdate(cf ConfigFile) ([]*Difference, error) {
	orbs, err := cf.Parse()
	if err != nil {
		return nil, err
	}

	cl := NewDefaultClient()
	differences := make([]*Difference, 0, len(orbs))
	for _, o := range orbs {
		if !o.Version().IsSemantic() {
			continue
		}

		newVersion, err := cl.LatestVersion(o)
		if err != nil {
			return nil, err
		}

		diff := NewDifference(o, newVersion)
		if !diff.HasUpdate() {
			continue
		}

		if o.Version() != newVersion.Version() {
			differences = append(differences, diff)
		}
	}

	return differences, nil
}
