package gopkg_test

import (
	"os"
	"reflect"
	"testing"

	. "github.com/goph/licensei/pkg/pkgutil/gopkg"
)

func TestReadLock(t *testing.T) {
	tests := map[string]*Lock{
		"testdata/lock/golden0.toml": {
			SolveMeta: SolveMeta{InputImports: []string{}},
			Projects: []Project{
				{
					Branch:    "master",
					Digest:    "1:666f6f",
					Name:      "github.com/golang/dep",
					Packages:  []string{"."},
					PruneOpts: "",
					Revision:  "d05d5aca9f895d19e9265839bffeadd74a2d2ecb",
				},
			},
		},

		"testdata/lock/golden1.toml": {
			SolveMeta: SolveMeta{InputImports: []string{}},
			Projects: []Project{
				{
					Version:   "0.12.2",
					Digest:    "1:666f6f",
					Name:      "github.com/golang/dep",
					Packages:  []string{"."},
					PruneOpts: "NUT",
					Revision:  "d05d5aca9f895d19e9265839bffeadd74a2d2ecb",
				},
			},
		},
	}

	for file, want := range tests {
		t.Run(file, func(t *testing.T) {
			lockFile, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}

			got, err := ReadLock(lockFile)
			if err != nil {
				t.Fatalf("should have read Lock correctly, but got err %q", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Error("valid lock did not parse as expected")
			}
		})
	}
}
