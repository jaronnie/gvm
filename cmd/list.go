/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "gvm list go bin in local or remote",
	Long:  `gvm list go bin in local or remote`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
