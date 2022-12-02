package licensei

import (
	"github.com/goph/licensei/pkg/pkgmgr/gomod"
	"golang.org/x/exp/slog"

	"github.com/pkg/errors"
)

type gomodDependencySource struct {
	logger *slog.Logger
}

func NewGoModDependencySource(logger *slog.Logger) *gomodDependencySource {
	return &gomodDependencySource{
		logger: logger,
	}
}

func (s *gomodDependencySource) Dependencies() ([]Dependency, error) {
	logger := s.logger.WithGroup("gomod")

	logger.Debug("listing go modules")

	deps, err := gomod.ListDeps("")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list dependencies")
	}

	var dependencies []Dependency // nolint: prealloc
	moduleMap := map[string]bool{}

	for _, dep := range deps {
		// Skip the main module and non-module dependencies
		if dep.Module == nil || dep.Module.Main {
			continue
		}

		// Skip already processed modules
		if _, ok := moduleMap[dep.Module.Path]; ok {
			continue
		}

		moduleMap[dep.Module.Path] = true

		dependencies = append(
			dependencies,
			Dependency{
				Name:     dep.Module.Path,
				Revision: dep.Module.Version,
				Type:     "gomod",
			},
		)
	}

	return dependencies, nil
}
