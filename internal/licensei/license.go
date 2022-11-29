package licensei

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	githubdetector "github.com/goph/licensei/pkg/detector/github"
	"github.com/goph/licensei/pkg/detector/sourced"
	"github.com/goph/licensei/pkg/licensematch"
	"github.com/goph/licensei/pkg/resolver"
)

type LicenseDetector struct {
	githubDetectorOptions []githubdetector.DetectorOption
}

func NewLicenseDetector(githubToken string) *LicenseDetector {
	l := new(LicenseDetector)

	if githubToken != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		l.githubDetectorOptions = append(l.githubDetectorOptions, githubdetector.Client(client))
	}

	return l
}

func (d *LicenseDetector) Detect(dependencies []Dependency) ([]Dependency, error) {
	var detector interface {
		Detect() (map[string]float32, error)
	}

	for key, dep := range dependencies {
		if dep.License != "" {
			continue
		}

		var matches map[string]float32

		_, err := os.Stat("vendor/" + dep.Name)
		if !os.IsNotExist(err) {
			f, err := sourced.FilerFromDirectory("vendor/" + dep.Name)
			if err != nil {
				return nil, emperror.With(
					errors.Wrap(err, "could not initialize license detector"),
					"dependency", dep.Name,
				)
			}
			detector = sourced.NewDetector(f)

			matches, err = detector.Detect()
			if err != nil { // TODO: add error handling
				continue
			}
		}

		name := resolver.Resolve(dep.Name)

		if strings.HasPrefix(name, "github.com") {
			repoData := strings.Split(name, "/")
			detector = githubdetector.NewDetector(repoData[1], repoData[2], d.githubDetectorOptions...)

			m, err := detector.Detect()
			if err == nil { // TODO: add error handling
				matches = licensematch.Merge(matches, m)
			}
		}

		license, confidence := licensematch.Best(matches)

		dependencies[key].License = license
		dependencies[key].Confidence = confidence
	}

	return dependencies, nil
}
