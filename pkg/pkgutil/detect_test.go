package pkgutil_test

import (
	"reflect"
	"testing"

	. "github.com/goph/licensei/pkg/pkgutil"
)

func TestDetectPackageManagers(t *testing.T) {
	tests := map[string]*PackageManagers{
		"dep": {
			Dep: true,
		},
		"gomod": {
			GoMod: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			pkgmgrs, err := DetectPackageManagers("testdata/detect/" + name)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(pkgmgrs, test) {
				t.Error("package manager detection failed")
			}
		})
	}
}
