/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/utilx"
)

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
	cobra.OnInitialize(initGlobalValue)
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if b, _ := utilx.PathExists(global.GVM_CONFIG_DIR); !b {
		if err := os.Mkdir(global.GVM_CONFIG_DIR, 0744); err != nil {
			panic(err)
		}
	}

	if b, _ := utilx.PathExists(global.GVM_CONFIG_FILE); !b {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(global.GVM_CONFIG_DIR)

		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(global.GVM_CONFIG_FILE)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read in config meet error. Err: [%v]", err)
	}
}

func initGlobalValue() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	global.GVM_CONFIG_DIR = fmt.Sprintf("%s/gvm", homeDir)
	global.GVM_CONFIG_FILE = fmt.Sprintf("%s/gvm/config.toml", homeDir)
	global.GVM_GOROOT = fmt.Sprintf("%s/gvm/goroot", homeDir)
}
