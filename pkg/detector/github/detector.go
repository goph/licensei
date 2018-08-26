package github

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type detector struct {
	owner string
	repo  string

	client *github.Client
}

func NewDetector(owner string, repo string) *detector {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	d := &detector{
		owner:  owner,
		repo:   repo,
		client: github.NewClient(tc),
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
