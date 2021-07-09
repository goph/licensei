package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// DetectorOption configures the detector.
type DetectorOption interface {
	apply(d *detector)
}

// DetectorOptionFunc makes a function a detector option when it's definition is compatible.
type DetectorOptionFunc func(d *detector)

func (fn DetectorOptionFunc) apply(d *detector) {
	fn(d)
}

// Client configures a Github client.
func Client(c *github.Client) DetectorOption {
	return DetectorOptionFunc(func(d *detector) {
		d.client = c
	})
}

type detector struct {
	owner string
	repo  string

	client *github.Client
}

func NewDetector(owner string, repo string, opts ...DetectorOption) *detector {
	d := &detector{
		owner: owner,
		repo:  repo,
	}

	for _, opt := range opts {
		opt.apply(d)
	}

	// Default Github client
	if d.client == nil {
		d.client = github.NewClient(http.DefaultClient)
	}

	return d
}

func (d *detector) Detect() (map[string]float32, error) {
	return d.DetectContext(context.Background())
}

func (d *detector) DetectContext(ctx context.Context) (map[string]float32, error) {
	lic, _, err := d.client.Repositories.License(ctx, d.owner, d.repo)
	if err != nil {
		return nil, err
	}

	if lic.GetLicense().GetSPDXID() != "" && lic.GetLicense().GetSPDXID() != "NOASSERTION" {
		return map[string]float32{lic.GetLicense().GetSPDXID(): 1}, nil
	}

	// There is a license, but it couldn't be detected.
	if lic.GetLicense().GetSPDXID() == "NOASSERTION" {
		matches, err := licensedb.Detect(&filerImpl{License: lic})
		if err != nil {
			return nil, err
		}

		m := make(map[string]float32, len(matches))

		for l, v := range matches {
			m[l] = v.Confidence
		}

		return m, nil
	}

	return nil, errors.New("no license found")
}

// filerImpl implements filer.Filer to return the license text directly
// from the github.RepositoryLicense structure.
// Copied from https://github.com/mitchellh/golicense/blob/dafaeff2016e81a739a2346cade4afd87b8e8647/license/github/detect.go#L49
type filerImpl struct {
	License *github.RepositoryLicense
}

func (f *filerImpl) ReadFile(name string) ([]byte, error) {
	if name != "LICENSE" {
		return nil, fmt.Errorf("unknown file: %s", name)
	}

	return base64.StdEncoding.DecodeString(f.License.GetContent())
}

func (f *filerImpl) ReadDir(dir string) ([]filer.File, error) {
	// We only support root
	if dir != "" {
		return nil, nil
	}

	return []filer.File{{Name: "LICENSE"}}, nil
}

func (f *filerImpl) Close() {}

func (f *filerImpl) PathsAreAlwaysSlash() bool {
	return true
}
