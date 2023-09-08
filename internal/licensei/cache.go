package licensei

import (
	"bytes"
	"io"
	"os"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type licenseCache struct {
	Dependencies []Dependency `toml:"dependencies"`
}

// ReadCache reads the cache from file.
func ReadCache(r io.Reader) ([]Dependency, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read byte stream")
	}

	cache := new(licenseCache)
	err = toml.Unmarshal(buf.Bytes(), cache)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the cache as TOML")
	}

	return cache.Dependencies, nil
}

// WriteCache writes the generated cache to file.
func WriteCache(w io.Writer, dependencies []Dependency) error {
	encoder := toml.NewEncoder(w).ArraysWithOneElementPerLine(true)

	return encoder.Encode(licenseCache{Dependencies: dependencies})
}

type cacheDependencySource struct {
	delegatedDependencySource dependencySource
	logger                    *slog.Logger
}

func NewCacheProjectSource(delegatedDependencySource dependencySource, logger *slog.Logger) *cacheDependencySource {
	return &cacheDependencySource{
		delegatedDependencySource: delegatedDependencySource,
		logger:                    logger,
	}
}

func (s *cacheDependencySource) Dependencies() ([]Dependency, error) {
	logger := s.logger.WithGroup("dep_cache")

	var cachedDependencies []Dependency

	logger.Debug("load cache file")

	cacheFile, err := os.Open(".licensei.cache")
	defer func() {
		_ = cacheFile.Close()
	}()

	if err == nil {
		p, err := ReadCache(cacheFile)
		if err != nil {
			return nil, errors.WithMessage(err, "could not parse cache")
		}

		cachedDependencies = p
	} else if !os.IsNotExist(err) {
		// cache could not be loaded
		// at least log it dude
		return nil, err
	}

	logger.Debug("load dependencies")

	sourceDependencies, err := s.delegatedDependencySource.Dependencies()
	if err != nil {
		return nil, err
	}

	logger.Debug("process dependencies")

	cachedDependencyIndex := indexDependencies(cachedDependencies)

	// We don't know the exact size upfront, we could only estimate a maximum size
	// nolint: prealloc
	var dependencies []Dependency

	for _, dep := range sourceDependencies {
		cacheIndex, ok := cachedDependencyIndex[dep.Name]

		if ok {
			// Same revision, so license and confidence information is valid
			if cachedDependencies[cacheIndex].Revision == dep.Revision {
				dep.License = cachedDependencies[cacheIndex].License
				dep.Confidence = cachedDependencies[cacheIndex].Confidence
			}
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}
