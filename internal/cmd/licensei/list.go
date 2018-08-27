package licensei

import (
	"errors"
	"os"

	"github.com/goph/licensei/internal/licensei"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listOptions struct {
	format string

	githubToken string
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
			options.githubToken = viper.GetString("github_token")

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

	source := licensei.NewCacheProjectSource(licensei.NewDepDependencySource())
	dependencies, err := source.Dependencies()
	if err != nil {
		return err
	}

	detector := licensei.NewLicenseDetector(options.githubToken)

	dependencies, err = detector.Detect(dependencies)
	if err != nil {
		return err
	}

	var viewModel licensei.ListViewModel

	for _, dep := range dependencies {
		viewModel.Dependencies = append(
			viewModel.Dependencies,
			licensei.ListDependencyItem{
				Name:       dep.Name,
				License:    dep.License,
				Confidence: dep.Confidence,
			},
		)
	}

	return view.Render(viewModel)
}
