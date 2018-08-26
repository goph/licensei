package licensei

import (
	"os"
	"strings"

	"github.com/goph/licensei/internal/licensei"
	"github.com/goph/licensei/pkg/detector/github"
	"github.com/goph/licensei/pkg/detector/sourced"
	"github.com/goph/licensei/pkg/licensematch"
	"github.com/spf13/cobra"
)

type cacheOptions struct {
	update bool
}

func NewCacheCommand() *cobra.Command {
	var options cacheOptions

	cmd := &cobra.Command{
		Use:   "cache [OPTIONS]",
		Short: "Cache licenses of dependencies in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCache(options)
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&options.update, "update", false, "Invalidate and update the entire cache")

	return cmd
}
func runCache(options cacheOptions) error {
	source := licensei.NewCacheProjectSource(licensei.NewDepProjectSource())
	projects, err := source.Projects()
	if err != nil {
		return err
	}

	var detector interface {
		Detect() (map[string]float32, error)
	}

	for key, project := range projects {
		if project.License != "" {
			continue
		}

		f, err := sourced.FilerFromDirectory("vendor/" + project.Name)
		if err != nil {
			panic(err)
		}
		detector = sourced.NewDetector(f)

		matches, err := detector.Detect()
		if err != nil {
			continue
		}

		if strings.HasPrefix(project.Name, "github.com") {
			repoData := strings.Split(project.Name, "/")
			detector = github.NewDetector(repoData[1], repoData[2])

			m, err := detector.Detect()
			if err == nil {
				matches = licensematch.Merge(matches, m)
			}
		}

		license, confidence := licensematch.Best(matches)

		projects[key].License = license
		projects[key].Confidence = confidence
	}

	cacheFile, err := os.Create(".licensei.cache")
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	licensei.WriteCache(cacheFile, licensei.LicenseCache{Projects: projects})

	return nil
}
