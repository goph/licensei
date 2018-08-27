package github

import (
	"context"
	"net/http"

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

	if lic.License.SPDXID != nil {
		return map[string]float32{*(lic.License.SPDXID): 1}, nil
	}

	return nil, errors.New("no license found")
}
