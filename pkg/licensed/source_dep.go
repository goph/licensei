package licensed

import (
	"bytes"
	"io"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type lock struct {
	Projects []lockedProject `toml:"projects"`
}

type lockedProject struct {
	Name     string `toml:"name"`
	Revision string `toml:"revision"`
}

func readLock(r io.Reader) (*lock, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read byte stream")
	}

	raw := lock{}
	err = toml.Unmarshal(buf.Bytes(), &raw)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse the lock as TOML")
	}

	return &raw, nil
}

func DepSource() ([]Package, error) {
	lockFile, err := os.Open("Gopkg.lock")
	if err != nil {
		return nil, err
	}

	locked, err := readLock(lockFile)

	var packages []Package

	for _, project := range locked.Projects {
		pkg := Package{
			Name:     project.Name,
			Revision: project.Revision,
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}
