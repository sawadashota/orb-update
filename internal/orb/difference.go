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

// DetectUpdateSet reads multi config file and differences set
func DetectUpdateSet(cfs []ConfigFile) ([]*Difference, error) {
	set := NewDifferenceSet()
	for _, cf := range cfs {
		diffs, err := DetectUpdate(cf)
		if err != nil {
			return nil, err
		}

		set.addMulti(diffs...)
	}

	return set.set, nil
}

// differenceSet is set of Difference
// set should not be multiply
type differenceSet struct {
	set []*Difference
}

// NewDifferenceSet .
func NewDifferenceSet() *differenceSet {
	return &differenceSet{set: make([]*Difference, 0)}
}

// addMulti Differences
func (ds *differenceSet) addMulti(diffs ...*Difference) {
	for _, diff := range diffs {
		ds.add(diff)
	}
}

// add Differences
func (ds *differenceSet) add(diff *Difference) {
	for _, d := range ds.set {
		if d.Old.IsSameOrb(diff.Old) {
			return
		}
	}
	ds.set = append(ds.set, diff)
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
