package licensei

import (
	"fmt"
	"strings"

	"github.com/goph/licensei/internal/licensei"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type checkOptions struct {
	approved []string
}

func NewCheckCommand() *cobra.Command {
	var options checkOptions

	cmd := &cobra.Command{
		Use:   "check [OPTIONS]",
		Short: "Check licenses of dependencies in the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.approved = viper.GetStringSlice("approved")

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

	source := licensei.NewCacheProjectSource(licensei.NewDepProjectSource())
	projects, err := source.Projects()
	if err != nil {
		return err
	}

	var violations []licensei.Project

	for _, project := range projects {
		var approved bool

		for _, license := range options.approved {
			approved = approved || strings.ToLower(license) == strings.ToLower(project.License)
		}

		if !approved {
			violations = append(violations, project)
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

		return nil
	}

	fmt.Println("No license violations! Good job!")

	return nil
}
