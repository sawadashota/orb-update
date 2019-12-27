package extraction_test

import (
	"os"
	"testing"

	"github.com/sawadashota/orb-update/internal/extraction"

	"github.com/sawadashota/orb-update/internal/orb"
)

func containOrb(t *testing.T, needle *orb.Orb, haystack []*orb.Orb) bool {
	t.Helper()

	for _, o := range haystack {
		if needle.String() == o.String() {
			return true
		}
	}

	return false
}

func TestExtraction_Orbs(t *testing.T) {
	cases := map[string]struct {
		configPath string
		want       []*orb.Orb
		wantErr    bool
	}{
		"correct format": {
			configPath: "./testdata/correct-format.yml",
			want: []*orb.Orb{
				orb.New("example", "example01", "3.4.1"),
				orb.New("example", "example02", "1.0.0"),
			},
			wantErr: false,
		},
		"no orb": {
			configPath: "./testdata/no-orb.yml",
			want:       []*orb.Orb{},
			wantErr:    false,
		},
		"incorrect format": {
			configPath: "./testdata/incorrect-format.yml",
			want:       nil,
			wantErr:    true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			file, err := os.Open(c.configPath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			cf, err := extraction.NewExtraction(file)
			if err != nil {
				t.Fatal(err)
			}

			got, err := cf.Orbs()
			if (err != nil) != c.wantErr {
				t.Errorf("Orbs() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if err != nil {
				return
			}

			if len(got) != len(c.want) {
				t.Errorf("Orbs() got = %v, want %v", got, c.want)
			}

			for _, gotOrb := range got {
				if !containOrb(t, gotOrb, c.want) {
					t.Errorf("Orbs() got = %v, want %v", got, c.want)
				}
			}
		})
	}
}
