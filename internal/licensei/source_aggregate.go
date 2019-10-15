package licensei

import (
	"github.com/goph/licensei/pkg/pkgmgr"
)

type aggregatedDependencySource struct {
	dependencySources []dependencySource
}

func NewAggregatedDependencySource() *aggregatedDependencySource {
	pkgmgrs, err := pkgmgr.DetectPackageManagers(".")
	if err != nil {
		panic(err)
	}

	source := &aggregatedDependencySource{
		dependencySources: []dependencySource{},
	}

	if pkgmgrs.Dep {
		source.dependencySources = append(source.dependencySources, NewDepDependencySource())
	}

	if pkgmgrs.GoMod {
		source.dependencySources = append(source.dependencySources, NewGoModDependencySource())
	}

	return source
}

func (s *aggregatedDependencySource) Dependencies() ([]Dependency, error) {
	var deps []Dependency // nolint: prealloc

	for _, depSource := range s.dependencySources {
		sdeps, err := depSource.Dependencies()
		if err != nil {
			return nil, err
		}

		deps = append(deps, sdeps...)
	}

	return deps, nil
}
