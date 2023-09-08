package licensei

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/pkg/pkgmgr/gomod"
)

type gomodDependencySource struct {
	logger *slog.Logger
	path   string
}

func NewGoModDependencySource(logger *slog.Logger, path string) *gomodDependencySource {
	return &gomodDependencySource{
		logger: logger,
		path:   path,
	}
}

func (s *gomodDependencySource) Dependencies() ([]Dependency, error) {
	logger := s.logger.WithGroup("gomod")

	logger.Debug("listing go modules")

	deps, err := gomod.ListDeps(s.path)
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
