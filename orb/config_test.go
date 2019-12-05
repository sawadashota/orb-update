package orb_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/sawadashota/orb-update/orb"
)

func TestConfigFile_Parse(t *testing.T) {
	cases := map[string]struct {
		configPath string
		want       *orb.Config
		wantErr    bool
	}{
		"correct format": {
			configPath: "./testdata/correct-format.yml",
			want: &orb.Config{
				Orbs: []*orb.Orb{
					orb.NewOrb("example", "example01", "3.4.1"),
					orb.NewOrb("example", "example02", "1.0.0"),
				},
			},
			wantErr: false,
		},
		"no orb": {
			configPath: "./testdata/no-orb.yml",
			want:       new(orb.Config),
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
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Parse() got = %v, want %v", got, c.want)
			}
		})
	}
}

func TestConfigFile_Update(t *testing.T) {
	type args struct {
		newVersions []*orb.Orb
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
				newVersions: []*orb.Orb{
					orb.NewOrb("example", "example01", "3.4.2"),
					orb.NewOrb("example", "example02", "1.1.0"),
				},
			},
			want: `orbs:
  example01: example/example01@3.4.2
  example02: example/example02@1.1.0

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
				newVersions: []*orb.Orb{
					orb.NewOrb("example", "example01", "3.4.2"),
					orb.NewOrb("example", "example02", "1.1.0"),
				},
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
				newVersions: []*orb.Orb{
					orb.NewOrb("example", "example01", "3.4.2"),
					orb.NewOrb("example", "example02", "1.1.0"),
				},
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
			err = cf.Update(&result, c.args.newVersions...)
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
