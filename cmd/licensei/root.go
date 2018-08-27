package main

import (
	"fmt"
	"os"

	"github.com/goph/licensei/internal/cmd/licensei"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file (default is $PWD/.licensei.yaml)")

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
}

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
