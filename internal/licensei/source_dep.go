package licensei

import (
	"os"

	"github.com/goph/licensei/pkg/pkgutil/gopkg"
)

type depProjectSource struct {
}

func NewDepProjectSource() *depProjectSource {
	return new(depProjectSource)
}

func (s *depProjectSource) Projects() ([]Project, error) {
	lockFile, err := os.Open("Gopkg.lock")
	if err != nil {
		return nil, err
	}

	locked, err := gopkg.ReadLock(lockFile)

	var packages []Project

	for _, project := range locked.Projects {
		pkg := Project{
			Name:     project.Name,
			Revision: project.Revision,
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}
