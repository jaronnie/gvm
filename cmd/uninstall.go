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
	"github.com/jaronnie/gvm/internal/vm"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "gvm uninstall",
	Long:  `gvm uninstall`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		rd := vm.NewLocalVM()

		vs, err := rd.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return vs, cobra.ShellCompDirectiveDefault
	},
	RunE: uninstall,
}

func uninstall(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	goRoot, _ := os.Readlink(global.GvmGoroot)
	if filepath.Join(global.GvmConfigDir, gov) == goRoot {
		return errors.Errorf("can not uninstall %s, please exec `gvm activate go<other_version>`", gov)
	}

	err := os.RemoveAll(filepath.Join(global.GvmConfigDir, gov))
	if err != nil {
		return err
	}
	fmt.Printf("🚀Uninstall %s successfully\n", gov)

	return nil
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
