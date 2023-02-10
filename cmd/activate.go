/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "gvm activate go",
	Long:  `gvm activate go`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:  activate,
}

func activate(cmd *cobra.Command, args []string) error {
	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	homeDir, _ := os.UserHomeDir()

	gvmConfigPath := fmt.Sprintf(GVMConfigPath, homeDir)

	goBaseRoot := filepath.Join(gvmConfigPath, "goroot")
	goRoot := filepath.Join(gvmConfigPath, gov)

	_ = os.Remove(goBaseRoot)

	err := os.Symlink(goRoot, goBaseRoot)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
