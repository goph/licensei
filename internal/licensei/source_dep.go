package licensei

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/goph/licensei/pkg/pkgmgr/gopkg"
)

type depDependencySource struct {
	path string
}

func NewDepDependencySource(path string) *depDependencySource {
	return &depDependencySource{
		path: path,
	}
}

func (s *depDependencySource) Dependencies() ([]Dependency, error) {
	lockFile, err := os.Open(filepath.Join(s.path, "Gopkg.lock"))
	if err != nil {
		return nil, errors.Wrap(err, "could not open dep lock file")
	}
	defer lockFile.Close()

	lock, err := gopkg.ReadLock(lockFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read dep lock file")
	}

	dependencies := make([]Dependency, len(lock.Projects))

	for i, project := range lock.Projects {
		pkg := Dependency{
			Name:     project.Name,
			Revision: project.Revision,
			Type:     "dep",
		}
		dependencies[i] = pkg
	}

	return dependencies, nil
}
