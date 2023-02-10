/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"github.com/jaronnie/gvm/utilx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const GVMConfigPath = "%s/.gvm"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "golang version manage",
	Long:  `golang version manage`,
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

	gvmConfigPath := fmt.Sprintf(GVMConfigPath, homeDir)

	if b, _ := utilx.PathExists(gvmConfigPath); !b {
		_ = os.Mkdir(gvmConfigPath, 0744)
	}

	cfgFile := gvmConfigPath + "/config.toml"

	if b, _ := utilx.PathExists(cfgFile); !b {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(gvmConfigPath)

		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(cfgFile)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read in config meet error. Err: [%v]", err)
	}
}
