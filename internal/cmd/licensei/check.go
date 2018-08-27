package licensei

import (
	"fmt"
	"os"
	"strings"

	"github.com/goph/licensei/internal/licensei"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type checkOptions struct {
	approved []string
	ignored  []string

	githubToken string
}

func NewCheckCommand() *cobra.Command {
	var options checkOptions

	cmd := &cobra.Command{
		Use:   "check [OPTIONS]",
		Short: "Check licenses of dependencies in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.approved = viper.GetStringSlice("approved")
			options.ignored = viper.GetStringSlice("ignored")

			options.githubToken = viper.GetString("github_token")

			return runCheck(options)
		},
	}

	return cmd
}
func runCheck(options checkOptions) error {
	if len(options.approved) == 0 {
		fmt.Println("everything is approved")

		return nil
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

	var violations []licensei.Dependency

	ignored := make(map[string]bool, len(options.ignored))

	for _, name := range options.ignored {
		ignored[name] = true
	}

	for _, dep := range dependencies {
		var approved bool

		for _, license := range options.approved {
			approved = approved || strings.ToLower(license) == strings.ToLower(dep.License)
		}

		if _, ignore := ignored[dep.Name]; !approved && !ignore {
			violations = append(violations, dep)
		}
	}

	if len(violations) > 0 {
		for _, project := range violations {
			if project.License == "" {
				fmt.Printf("%s: no license file found\n", project.Name)
			} else {
				fmt.Printf("%s: license violation: %s\n", project.Name, project.License)
			}
		}

		os.Exit(2)

		return nil
	}

	fmt.Println("No license violations! Good job!")

	return nil
}
