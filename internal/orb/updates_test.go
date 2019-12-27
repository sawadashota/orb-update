package orb_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/sawadashota/orb-update/internal/orb"
)

func init() {
	orb.CircleCIClient = new(testClient)
}

type testClient struct{}

func (tc *testClient) LatestVersion(o *orb.Orb) (*orb.Orb, error) {
	latestOrbs := map[string]*orb.Orb{
		"example/example01": orb.New("example", "example01", "3.4.2"),
		"example/example02": orb.New("example", "example02", "1.0.1"),
		"example/example03": orb.New("example", "example03", "5.0.9"),
	}

	packageName := fmt.Sprintf("%s/%s", o.Namespace(), o.Name())
	found, ok := latestOrbs[packageName]
	if !ok {
		return nil, errors.Errorf("no such a orb: %s", packageName)
	}
	return found, nil
}

func TestExtraction_Updates(t *testing.T) {
	cases := map[string]struct {
		filters    []orb.Filter
		configBody string
		want       []*orb.Update
		wantErr    bool
	}{
		"duplicated orbs at config": {
			filters: orb.ExcludeMatchPackages([]string{"example/example03"}),
			configBody: `
orbs:
  example01: example/example01@3.4.1
  example01: example/example01@3.4.1
  example02: example/example02@1.0.0
  example03: example/example03@2.0.0

job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"`,
			want: []*orb.Update{
				orb.NewUpdate(
					orb.New("example", "example01", "3.4.1"),
					orb.New("example", "example01", "3.4.2"),
				),
				orb.NewUpdate(
					orb.New("example", "example02", "1.0.0"),
					orb.New("example", "example02", "1.0.1"),
				),
			},
			wantErr: false,
		},
		"all orbs are updated": {
			filters: orb.ExcludeMatchPackages([]string{"example/example03"}),
			configBody: `
orbs:
  example01: example/example01@3.4.2
  example02: example/example02@1.0.1
  example03: example/example03@2.0.0

job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"`,
			want:    []*orb.Update{},
			wantErr: false,
		},
		"un-exists orb": {
			filters: orb.ExcludeMatchPackages([]string{"example/example03"}),
			configBody: `
orbs:
  example01: example/example01@3.4.1
  example02: example/example02@1.0.0
  example03: example/example03@2.0.0
  example04: example/example04@8.0.0

job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"`,
			want:    nil,
			wantErr: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			b := bytes.NewBufferString(c.configBody)
			e, err := orb.NewExtraction(b)
			if err != nil {
				t.Fatal(err)
			}
			got, err := e.Updates(c.filters...)
			if (err != nil) != c.wantErr {
				t.Errorf("Updates() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			for _, gotUpdate := range got {
				if !containUpdate(t, gotUpdate, c.want) {
					t.Errorf("Updates() got = %v, want %v", got, c.want)
				}
			}
		})
	}
}

func containUpdate(t *testing.T, needle *orb.Update, haystack []*orb.Update) bool {
	t.Helper()

	for _, update := range haystack {
		if needle.Before.String() == update.Before.String() && needle.After.String() == update.After.String() {
			return true
		}
	}

	return false
}
