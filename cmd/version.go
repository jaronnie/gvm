/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get version",
	Long:  `get version`,
	RunE:  version,
}

func version(cmd *cobra.Command, args []string) error {
	fmt.Println(Version)
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
