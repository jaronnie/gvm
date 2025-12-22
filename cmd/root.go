/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/utilx"
)

var (
	Registry string
	Cache    bool
)

const GVM = "gvm"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "golang version manage",
	Long:  `golang version manage`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// isInitialized checks if GVM has been initialized by looking for the .gvmrc file
func isInitialized() bool {
	if global.GvmConfigRc == "" {
		return false
	}
	exists, _ := utilx.PathExists(global.GvmConfigRc)
	return exists
}

// checkInit validates that GVM has been initialized before allowing commands to run
func checkInit(cmd *cobra.Command, args []string) error {
	// Allow init command to run without checking initialization
	if cmd.Name() == "init" || cmd.Name() == "version" || cmd.Name() == "help" {
		return nil
	}

	if !isInitialized() {
		return errors.New("gvm has not been initialized. Please run 'gvm init' first")
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Add persistent pre-run hook to check initialization for all commands
	rootCmd.PersistentPreRunE = checkInit

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initGlobalValue)
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&Registry, "registry", "i", "https://dl.google.com/go", "set registry")
	rootCmd.PersistentFlags().BoolVarP(&Cache, "cache", "", false, "get remote version by cache")

	registerFlagCompletion()
}

func initConfig() {
	if b, _ := utilx.PathExists(global.GvmConfigDir); !b {
		if err := os.Mkdir(global.GvmConfigDir, 0o744); err != nil {
			panic(err)
		}
	}

	if b, _ := utilx.PathExists(global.GvmConfigFile); !b {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(global.GvmConfigDir)

		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(global.GvmConfigFile)
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

	// check
	if stat, err := os.Stat(filepath.Join(homeDir, rootCmd.Use)); err == nil {
		if !stat.IsDir() {
			panic("please make sure $HOME/gvm is a dir")
		}
	}

	global.HomeDir = homeDir
	global.GvmConfigDir = fmt.Sprintf("%s/%s", homeDir, GVM)
	global.GvmConfigFile = fmt.Sprintf("%s/%s/config.toml", homeDir, GVM)
	global.GvmConfigRc = fmt.Sprintf("%s/%s/.gvmrc", homeDir, GVM)
	global.GvmGoroot = fmt.Sprintf("%s/%s/goroot", homeDir, GVM)
}

func registerFlagCompletion() {
	_ = configCmd.RegisterFlagCompletionFunc("registry", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cobra.CompDebugln(strings.Join(args, ","), true)
		return SupportMirrors, cobra.ShellCompDirectiveDefault
	})
}
