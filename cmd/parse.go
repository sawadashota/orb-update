package cmd

import (
	"github.com/sawadashota/orb-update/orb"
)

func parse(cf *orb.ConfigFile, filePath string) ([]*orb.Difference, error) {
	conf, err := cf.Parse()
	if err != nil {
		return nil, err
	}

	cl := orb.NewDefaultClient()
	differences := make([]*orb.Difference, 0, len(conf.Orbs))
	for _, o := range conf.Orbs {
		newVersion, err := cl.LatestVersion(o)
		if err != nil {
			return nil, err
		}

		if o.Version() != newVersion.Version() {
			differences = append(differences, orb.NewDifference(o, newVersion))
		}
	}

	return differences, nil
}
