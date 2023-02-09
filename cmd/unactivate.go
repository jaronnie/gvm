/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// unactivateCmd represents the unactivate command
var unactivateCmd = &cobra.Command{
	Use:   "unactivate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unactivate called")
	},
}

func init() {
	rootCmd.AddCommand(unactivateCmd)
}
