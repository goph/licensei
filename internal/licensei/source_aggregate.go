package licensei

import (
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/pkg/pkgmgr"
)

type aggregatedDependencySource struct {
	dependencySources []dependencySource
}

func NewAggregatedDependencySource(logger *slog.Logger, path string) *aggregatedDependencySource {
	pkgmgrs, err := pkgmgr.DetectPackageManagers(path)
	if err != nil {
		panic(err)
	}

	source := &aggregatedDependencySource{
		dependencySources: []dependencySource{},
	}

	if pkgmgrs.Dep {
		source.dependencySources = append(source.dependencySources, NewDepDependencySource(path))
	}

	if pkgmgrs.GoMod {
		source.dependencySources = append(source.dependencySources, NewGoModDependencySource(logger, path))
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
