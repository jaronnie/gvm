/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install go bin",
	Long:  `install go bin`,
	RunE:  install,
}

func install(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)

}
