package orb

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

func (d *Difference) HasUpdate() bool {
	return d.Old.Version() != d.New.Version()
}