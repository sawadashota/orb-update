package cmd

import (
	"github.com/sawadashota/orb-update/configfile"
	"github.com/sawadashota/orb-update/orb"
)

func differences(cf *configfile.ConfigFile) ([]*orb.Difference, error) {
	orbs, err := cf.Parse()
	if err != nil {
		return nil, err
	}

	cl := orb.NewDefaultClient()
	differences := make([]*orb.Difference, 0, len(orbs))
	for _, o := range orbs {
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
