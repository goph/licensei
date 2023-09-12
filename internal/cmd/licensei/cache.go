package licensei

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/internal/licensei"
)

type cacheOptions struct {
	update bool

	githubToken string
}

func NewCacheCommand(globalOptions *GlobalOptions) *cobra.Command {
	var options cacheOptions

	cmd := &cobra.Command{
		Use:   "cache [OPTIONS]",
		Short: "Cache licenses of dependencies in the project",
		RunE: func(_ *cobra.Command, _ []string) error {
			options.githubToken = viper.GetString("github_token")

			return runCache(globalOptions, options)
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&options.update, "update", false, "Invalidate and update the entire cache")

	return cmd
}

func runCache(globalOptions *GlobalOptions, options cacheOptions) error {
	logger := slog.Default()

	logger.Debug("start cache")

	var dependencies []licensei.Dependency
	var err error

	// Invalidate cache data
	if options.update {
		source := licensei.NewAggregatedDependencySource(logger, globalOptions.Path)
		dependencies, err = source.Dependencies()
	} else {
		source := licensei.NewCacheProjectSource(
			licensei.NewAggregatedDependencySource(logger, globalOptions.Path),
			logger,
		)
		dependencies, err = source.Dependencies()
	}
	if err != nil {
		return err
	}

	detector := licensei.NewLicenseDetector(options.githubToken, logger)

	dependencies, err = detector.Detect(dependencies)
	if err != nil {
		return err
	}

	cacheFile, err := os.Create(".licensei.cache")
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	licensei.WriteCache(cacheFile, dependencies)

	return nil
}
