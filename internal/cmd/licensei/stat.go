package licensei

import (
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/internal/licensei"
)

type statOptions struct {
	githubToken string
}

func NewStatCommand(globalOptions *GlobalOptions) *cobra.Command {
	var options statOptions

	cmd := &cobra.Command{
		Use:   "stat [OPTIONS]",
		Short: "Create statistics about the licenses in the project",
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.githubToken = viper.GetString("github_token")

			return runStat(globalOptions, options, cmd.OutOrStderr())
		},
	}

	return cmd
}

func runStat(globalOptions *GlobalOptions, options statOptions, stdout io.Writer) error {
	logger := slog.Default()

	logger.Debug("start stat")

	source := licensei.NewCacheProjectSource(
		licensei.NewAggregatedDependencySource(logger, globalOptions.Path),
		logger,
	)
	dependencies, err := source.Dependencies()
	if err != nil {
		return err
	}

	detector := licensei.NewLicenseDetector(options.githubToken, logger)

	dependencies, err = detector.Detect(dependencies)
	if err != nil {
		return err
	}

	stats := map[string]int{}

	for _, dependency := range dependencies {
		stats[dependency.License] = stats[dependency.License] + 1
	}

	statKeys := make([]string, 0, len(stats))

	for lic := range stats {
		statKeys = append(statKeys, lic)
	}

	sort.Strings(statKeys)

	table := tablewriter.NewWriter(stdout)
	table.SetHeader([]string{"License", "Count"})

	for _, key := range statKeys {
		table.Append([]string{key, fmt.Sprintf("%d", stats[key])})
	}

	table.Render()

	return nil
}
