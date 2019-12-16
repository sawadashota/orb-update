package orb_test

import (
	"testing"

	"github.com/sawadashota/orb-update/orb"
)

func TestParseOrb(t *testing.T) {
	cases := map[string]struct {
		orb     string
		wantErr bool
	}{
		"correct format": {
			orb:     "example/example@1.1.11",
			wantErr: false,
		},
		"incorrect format": {
			orb:     "example@1.1.11",
			wantErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := orb.ParseOrb(c.orb)
			if c.wantErr && err == nil {
				t.Errorf("parse should be error but pass")
				return
			}

			if !c.wantErr && err != nil {
				t.Errorf("parse should be pass but error")
				return
			}
		})
	}
}

func TestVersion_IsSemantic(t *testing.T) {
	cases := map[string]struct {
		version string
		want    bool
	}{
		"major version": {
			version: "1.0.0",
			want:    true,
		},
		"beta version": {
			version: "1.0.0-beta",
			want:    true,
		},
		"non semantic version": {
			version: "hoge",
			want:    false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			version := orb.Version(c.version)
			if got := version.IsSemantic(); got != c.want {
				t.Errorf("IsSemantic() = %v, want %v", got, c.want)
			}
		})
	}
}
