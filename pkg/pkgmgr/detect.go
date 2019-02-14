package pkgmgr

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type PackageManagers struct {
	Dep   bool
	GoMod bool
}

func (p *PackageManagers) String() string {
	var pkgmgrs []string

	if p.Dep {
		pkgmgrs = append(pkgmgrs, "dep")
	}

	if p.GoMod {
		pkgmgrs = append(pkgmgrs, "gomod")
	}

	return strings.Join(pkgmgrs, ", ")
}

// DetectPackageManagers detects what kind of packages are used in a project.
func DetectPackageManagers(path string) (*PackageManagers, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var pkgmgrs PackageManagers

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Dep
		if file.Name() == "Gopkg.lock" {
			pkgmgrs.Dep = true
		}

		// Go modules
		if file.Name() == "go.mod" {
			pkgmgrs.GoMod = true
		}
	}

	return &pkgmgrs, nil
}
