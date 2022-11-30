package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	"github.com/goph/licensei/internal/cmd/licensei"
)

// nolint: gochecknoglobals
var config string

// nolint: gochecknoglobals
var debug bool

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file (default is $PWD/.licensei.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

	licensei.AddCommands(rootCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".licensei")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("can't read config:", err)
		os.Exit(1)
	}

	logHandler := slog.HandlerOptions{
		AddSource: debug,
		Level:     slog.InfoLevel,
	}

	if debug {
		logHandler.Level = slog.DebugLevel
	}

	logger := slog.New(logHandler.NewTextHandler(os.Stderr))

	slog.SetDefault(logger)
}

// nolint: gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "licensei",
	Short: "License master",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
