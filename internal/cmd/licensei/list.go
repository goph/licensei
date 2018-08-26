package licensei

import (
	"errors"
	"os"
	"strings"

	"github.com/goph/licensei/internal/licensei"
	"github.com/goph/licensei/pkg/detector/github"
	"github.com/goph/licensei/pkg/detector/sourced"
	"github.com/goph/licensei/pkg/licensematch"
	"github.com/spf13/cobra"
)

type listOptions struct {
	format string
}

type listView interface {
	Render(model licensei.ListViewModel) error
}

func NewListCommand() *cobra.Command {
	var options listOptions

	cmd := &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "List licenses of dependencies in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.format, "format", "table", "Output format (table, json)")

	return cmd
}

func runList(options listOptions) error {
	var view listView
	switch options.format {
	case "", "table":
		view = licensei.NewTableListView(os.Stdout)

	case "json":
		view = licensei.NewJsonListView(os.Stdout)
	default:
		return errors.New("unsupported format: " + options.format)
	}

	source := licensei.NewCacheProjectSource(licensei.NewDepProjectSource())
	projects, err := source.Projects()
	if err != nil {
		return err
	}

	var detector interface {
		Detect() (map[string]float32, error)
	}

	var viewModel licensei.ListViewModel

	for _, project := range projects {
		if project.License != "" {
			viewModel.Projects = append(
				viewModel.Projects,
				licensei.ListProjectItem{
					Name:       project.Name,
					License:    project.License,
					Confidence: project.Confidence,
				},
			)

			continue
		}

		f, err := sourced.FilerFromDirectory("vendor/" + project.Name)
		if err != nil {
			panic(err)
		}
		detector = sourced.NewDetector(f)

		matches, err := detector.Detect()
		if err != nil {
			viewModel.Projects = append(
				viewModel.Projects,
				licensei.ListProjectItem{
					Name: project.Name,
				},
			)

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

		viewModel.Projects = append(
			viewModel.Projects,
			licensei.ListProjectItem{
				Name:       project.Name,
				License:    license,
				Confidence: confidence,
			},
		)
	}

	return view.Render(viewModel)
}
