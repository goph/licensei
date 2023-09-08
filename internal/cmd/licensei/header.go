package licensei

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goph/licensei/internal/licensei"
)

type headerOptions struct {
	template    string
	ignorePaths []string
	ignoreFiles []string
	authors     []string
}

func NewHeaderCommand(globalOptions *GlobalOptions) *cobra.Command {
	var options headerOptions

	cmd := &cobra.Command{
		Use:   "header",
		Short: "Check license header of files",
		RunE: func(_ *cobra.Command, _ []string) error {
			options.template = viper.GetString("header.template")
			options.ignorePaths = viper.GetStringSlice("header.ignorePaths")
			options.ignoreFiles = viper.GetStringSlice("header.ignoreFiles")
			options.authors = viper.GetStringSlice("header.authors")

			return runHeader(globalOptions, options)
		},
	}

	return cmd
}

func runHeader(globalOptions *GlobalOptions, options headerOptions) (err error) {
	target := globalOptions.Path

	if target == "" {
		target, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	violations, err := licensei.HeaderChecker{
		IgnorePaths: options.ignorePaths,
		IgnoreFiles: options.ignoreFiles,
		Authors:     options.authors,
	}.Check(target, options.template)
	if err != nil {
		return err
	}

	for path, violation := range violations {
		fmt.Printf("%s: %s\n", violation, path)
	}

	if len(violations) > 0 {
		os.Exit(1)
	}

	return nil
}
