package licensei

import (
	"os"

	"github.com/goph/licensei/pkg/pkgutil/gopkg"
	"github.com/pkg/errors"
)

type depDependencySource struct {
}

func NewDepDependencySource() *depDependencySource {
	return new(depDependencySource)
}

func (s *depDependencySource) Dependencies() ([]Dependency, error) {
	lockFile, err := os.Open("Gopkg.lock")
	if err != nil {
		return nil, errors.Wrap(err, "could not open dep lock file")
	}

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
