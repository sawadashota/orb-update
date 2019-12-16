package cmd

import (
	"github.com/sawadashota/orb-update/orb"
)

func differences(cf *orb.ConfigFile) ([]*orb.Difference, error) {
	conf, err := cf.Parse()
	if err != nil {
		return nil, err
	}

	cl := orb.NewDefaultClient()
	differences := make([]*orb.Difference, 0, len(conf.Orbs))
	for _, o := range conf.Orbs {
		if !o.Version().IsSemantic() {
			continue
		}

		newVersion, err := cl.LatestVersion(o)
		if err != nil {
			return nil, err
		}

		diff := orb.NewDifference(o, newVersion)
		if !diff.HasUpdate() {
			continue
		}

		if o.Version() != newVersion.Version() {
			differences = append(differences, diff)
		}
	}

	return differences, nil
}
