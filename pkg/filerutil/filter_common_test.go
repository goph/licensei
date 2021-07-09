package filerutil_test

import (
	"testing"

	"github.com/go-enry/go-license-detector/v4/licensedb/filer"

	. "github.com/goph/licensei/pkg/filerutil"
)

func TestCorrectnessFilter(t *testing.T) {
	tests := map[string]struct {
		file   filer.File
		filter bool
	}{
		"exclude": {
			file: filer.File{
				Name:  "license.docs",
				IsDir: false,
			},
			filter: false,
		},
		"include": {
			file: filer.File{
				Name:  "license",
				IsDir: false,
			},
			filter: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			got := CorrectnessFilter(test.file)

			if got != test.filter {
				t.Errorf("expected the filter value to be '%v', got '%v'", test.filter, got)
			}
		})
	}
}

func TestDirFilter(t *testing.T) {
	tests := map[string]struct {
		file   filer.File
		filter bool
	}{
		"nondir": {
			file: filer.File{
				Name:  "license",
				IsDir: false,
			},
			filter: true,
		},
		"dir": {
			file: filer.File{
				Name:  "cmd",
				IsDir: true,
			},
			filter: false,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			got := DirFilter(test.file)

			if got != test.filter {
				t.Errorf("expected the filter value to be '%v', got '%v'", test.filter, got)
			}
		})
	}
}
