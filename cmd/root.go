/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"github.com/jaronnie/gvm/utilx"
	"github.com/spf13/cobra"
	"os"
)

const GVMConfigPath = "%s/.gvm"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if b, _ := utilx.PathExists(fmt.Sprintf(GVMConfigPath, homeDir)); !b {
		_ = os.Mkdir(fmt.Sprintf(GVMConfigPath, homeDir), 0744)
	}
}
