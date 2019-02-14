package licensei

import "github.com/goph/licensei/pkg/pkgutil"

type aggregatedDependencySource struct {
	dependencySources []dependencySource
}

func NewAggregatedDependencySource() *aggregatedDependencySource {
	pkgmgrs, err := pkgutil.DetectPackageManagers(".")
	if err != nil {
		panic(err)
	}

	source := &aggregatedDependencySource{
		dependencySources: []dependencySource{},
	}

	if pkgmgrs.Dep {
		source.dependencySources = append(source.dependencySources, NewDepDependencySource())
	}

	return source
}

func (s *aggregatedDependencySource) Dependencies() ([]Dependency, error) {
	var deps []Dependency

	for _, depSource := range s.dependencySources {
		sdeps, err := depSource.Dependencies()
		if err != nil {
			return nil, err
		}

		for _, dep := range sdeps {
			deps = append(deps, dep)
		}
	}

	return deps, nil
}
