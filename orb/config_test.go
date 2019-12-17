package orb_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/sawadashota/orb-update/orb"
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

func TestConfigFile_Parse(t *testing.T) {
	cases := map[string]struct {
		configPath string
		want       []*orb.Orb
		wantErr    bool
	}{
		"correct format": {
			configPath: "./testdata/correct-format.yml",
			want: []*orb.Orb{
				orb.NewOrb("example", "example01", "3.4.1"),
				orb.NewOrb("example", "example02", "1.0.0"),
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

			cf, err := orb.NewConfigFile(file)
			if err != nil {
				t.Fatal(err)
			}

			got, err := cf.Parse()
			if (err != nil) != c.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if err != nil {
				return
			}

			if len(got) != len(c.want) {
				t.Errorf("Parse() got = %v, want %v", got, c.want)
			}

			for _, gotOrb := range got {
				if !containOrb(t, gotOrb, c.want) {
					t.Errorf("Parse() got = %v, want %v", got, c.want)
				}
			}
		})
	}
}

func TestConfigFile_Update(t *testing.T) {
	type args struct {
		diff *orb.Difference
	}
	cases := map[string]struct {
		configPath string
		args       args
		want       string
		wantErr    bool
	}{
		"correct format": {
			configPath: "./testdata/correct-format.yml",
			args: args{
				diff: orb.NewDifference(
					orb.NewOrb("example", "example01", "3.4.1"),
					orb.NewOrb("example", "example01", "3.4.2"),
				),
			},
			want: `orbs:
  example01: example/example01@3.4.2
  example02: example/example02@1.0.0

job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"
`,
			wantErr: false,
		},
		"no orb": {
			configPath: "./testdata/no-orb.yml",
			args: args{
				diff: orb.NewDifference(
					orb.NewOrb("example", "example01", "3.4.1"),
					orb.NewOrb("example", "example01", "3.4.2"),
				),
			},
			want: `job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"
`,
			wantErr: false,
		},
		"incorrect format": {
			configPath: "./testdata/incorrect-format.yml",
			args: args{
				diff: orb.NewDifference(
					orb.NewOrb("example", "example01", "3.4.1"),
					orb.NewOrb("example", "example01", "3.4.2"),
				),
			},
			want: `orbs:
  example01: example@3.4.1
  example02: example

job:
  example:
    docker:
      - image: example
    steps:
      - run: echo "Hello World"
`,
			wantErr: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			file, err := os.Open(c.configPath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			cf, err := orb.NewConfigFile(file)
			if err != nil {
				t.Fatal(err)
			}

			var result bytes.Buffer
			err = cf.Update(&result, c.args.diff)
			if (err != nil) != c.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if result.String() != c.want {
				t.Errorf("Update() \n===got===\n%s\n===want===\n%s", result.String(), c.want)
			}
		})
	}
}
