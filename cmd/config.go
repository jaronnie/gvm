/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SupportMirrors = []string{
	"https://dl.google.com/go",
	"https://mirrors.aliyun.com/golang",
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "gvm config",
	Long:  `gvm config`,
	RunE:  config,
}

func config(cmd *cobra.Command, args []string) error {
	viper.Set("registry", Registry)
	return viper.WriteConfig()
}

func init() {
	rootCmd.AddCommand(configCmd)
}
