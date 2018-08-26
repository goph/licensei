package gopkg

import (
	"bytes"
	"io"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// LockName is the lock file name used by dep.
const LockName = "Gopkg.lock"

// Lock represents the raw lock file structure of dep.
type Lock struct {
	SolveMeta SolveMeta `toml:"solve-meta"`
	Projects  []Project `toml:"projects"`
}

// SolveMeta holds metadata about the solving process that created the lock that
// is not specific to any individual project.
type SolveMeta struct {
	AnalyzerName    string   `toml:"analyzer-name"`
	AnalyzerVersion int      `toml:"analyzer-version"`
	SolverName      string   `toml:"solver-name"`
	SolverVersion   int      `toml:"solver-version"`
	InputImports    []string `toml:"input-imports"`
}

// Project represents a single project as seen in the dep lock file.
type Project struct {
	Name      string   `toml:"name"`
	Branch    string   `toml:"branch,omitempty"`
	Revision  string   `toml:"revision"`
	Version   string   `toml:"version,omitempty"`
	Source    string   `toml:"source,omitempty"`
	Packages  []string `toml:"packages"`
	PruneOpts string   `toml:"pruneopts"`
	Digest    string   `toml:"digest"`
}

// ReadLock parses data read as a lock file.
func ReadLock(r io.Reader) (*Lock, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read byte stream")
	}

	lock := new(Lock)
	err = toml.Unmarshal(buf.Bytes(), lock)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the lock as TOML")
	}

	return lock, nil
}
