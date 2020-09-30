package licensei

import (
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goph/licensei/internal/licensei"
)

type statOptions struct {
	githubToken string
}

func NewStatCommand() *cobra.Command {
	var options statOptions

	cmd := &cobra.Command{
		Use:   "stat [OPTIONS]",
		Short: "Create statistics about the licenses in the project",
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.githubToken = viper.GetString("github_token")

			return runStat(options, cmd.OutOrStderr())
		},
	}

	return cmd
}
func runStat(options statOptions, stdout io.Writer) error {
	source := licensei.NewCacheProjectSource(licensei.NewAggregatedDependencySource())
	dependencies, err := source.Dependencies()
	if err != nil {
		return err
	}

	detector := licensei.NewLicenseDetector(options.githubToken)

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
