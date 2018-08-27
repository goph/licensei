package licensei

import (
	"os"

	"github.com/goph/licensei/internal/licensei"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cacheOptions struct {
	update bool

	githubToken string
}

func NewCacheCommand() *cobra.Command {
	var options cacheOptions

	cmd := &cobra.Command{
		Use:   "cache [OPTIONS]",
		Short: "Cache licenses of dependencies in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.githubToken = viper.GetString("github_token")

			return runCache(options)
		},
	}

	flags := cmd.Flags()

	flags.BoolVar(&options.update, "update", false, "Invalidate and update the entire cache")

	return cmd
}
func runCache(options cacheOptions) error {
	var dependencies []licensei.Dependency
	var err error

	// Invalidate cache data
	if options.update {
		source := licensei.NewDepDependencySource()
		dependencies, err = source.Dependencies()
	} else {
		source := licensei.NewCacheProjectSource(licensei.NewDepDependencySource())
		dependencies, err = source.Dependencies()
	}
	if err != nil {
		return err
	}

	detector := licensei.NewLicenseDetector(options.githubToken)

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
