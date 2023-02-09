/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Registry string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "gvm config",
	Long:  `gvm config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&Registry, "registry", "i", "https://dl.google.com", "set registry")
}
