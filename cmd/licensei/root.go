package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/internal/cmd/licensei"
)

func newRootCommand(options *licensei.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licensei",
		Short: "License master",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initConfig(options)
		},
	}

	cmd.PersistentFlags().StringVar(&options.Config, "config", "", "config file (default is $PWD/.licensei.yaml)")
	cmd.PersistentFlags().StringVar(&options.Path, "path", "", "path to project (the current directory is used by default)")
	cmd.PersistentFlags().BoolVar(&options.Debug, "debug", false, "enable debug logging")

	return cmd
}

func initConfig(options *licensei.GlobalOptions) error {
	viper.AutomaticEnv()

	if options.Config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(options.Config)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".licensei")
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	logHandler := slog.HandlerOptions{
		AddSource: options.Debug,
		Level:     slog.InfoLevel,
	}

	if options.Debug {
		logHandler.Level = slog.DebugLevel
	}

	logger := slog.New(logHandler.NewTextHandler(os.Stderr))

	slog.SetDefault(logger)

	return nil
}

func Execute() {
	globalOptions := &licensei.GlobalOptions{}
	rootCmd := newRootCommand(globalOptions)
	licensei.AddCommands(rootCmd, globalOptions)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
