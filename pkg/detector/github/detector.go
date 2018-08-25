package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
)

type detector struct {
	owner string
	repo  string

	client *github.Client
}

func NewDetector(owner string, repo string) *detector {
	d := &detector{
		owner:  owner,
		repo:   repo,
		client: github.NewClient(http.DefaultClient),
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

	return map[string]float32{*lic.License.SPDXID: 1}, nil
}
