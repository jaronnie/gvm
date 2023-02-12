/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "gvm uninstall",
	Long:  `gvm uninstall`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:  uninstall,
}

func uninstall(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	goRoot, _ := os.Readlink(global.GVM_GOROOT)
	if filepath.Join(global.GVM_CONFIG_DIR, gov) == goRoot {
		return errors.Errorf("can not uninstall %s, please exec `gvm activate go<other_version>`", gov)
	}

	err := os.RemoveAll(filepath.Join(global.GVM_CONFIG_DIR, gov))
	if err != nil {
		return err
	}
	fmt.Printf("🚀Uninstall %s successfully\n", gov)

	return nil
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
