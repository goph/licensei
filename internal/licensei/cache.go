package licensei

import (
	"bytes"
	"io"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type LicenseCache struct {
	Projects []Project `toml:"projects"`
}

// ReadCache reads the cache from file.
func ReadCache(r io.Reader) (*LicenseCache, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read byte stream")
	}

	cache := new(LicenseCache)
	err = toml.Unmarshal(buf.Bytes(), cache)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the cache as TOML")
	}

	return cache, nil
}

// WriteCache writes the generated cache to file.
func WriteCache(w io.Writer, cache LicenseCache) error {
	encoder := toml.NewEncoder(w).ArraysWithOneElementPerLine(true)

	return encoder.Encode(cache)
}

type cacheProjectSource struct {
	delegatedProjectSource projectSource
}

func NewCacheProjectSource(delegatedProjectSource projectSource) *cacheProjectSource {
	return &cacheProjectSource{
		delegatedProjectSource: delegatedProjectSource,
	}
}

func (s *cacheProjectSource) Projects() ([]Project, error) {
	var projects []Project

	cacheFile, err := os.Open(".licensei.cache")
	if err == nil {
		cache, err := ReadCache(cacheFile)
		if err != nil {
			return nil, err
		}

		projects = cache.Projects
	} else if !os.IsNotExist(err) {
		// cache could not be loaded
		// at least log it dude
		return nil, err
	}

	cacheFile.Close()

	projectIndex := make(map[string]int, len(projects))

	for key, project := range projects {
		projectIndex[project.Name] = key
	}

	pkgs, err := s.delegatedProjectSource.Projects()
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		projectKey, ok := projectIndex[pkg.Name]
		if !ok {
			projects = append(
				projects,
				Project{
					Name:     pkg.Name,
					Revision: pkg.Revision,
				},
			)

			projectKey = len(projects) - 1
			projectIndex[pkg.Name] = projectKey
		}
	}

	return projects, nil
}
