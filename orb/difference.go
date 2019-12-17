package orb

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
