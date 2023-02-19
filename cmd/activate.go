/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/internal/vm"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "gvm activate go",
	Long:  `gvm activate go`,
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
	RunE: activate,
}

func activate(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	rd := vm.NewLocalVM()

	vs, err := rd.List()
	if err != nil {
		return err
	}

	if !lo.Contains(vs, gov) {
		return errors.Errorf("please exec `gvm install %s` first", gov)
	}

	goRoot := filepath.Join(global.GvmConfigDir, gov)

	_ = os.Remove(global.GvmGoroot)

	err = os.Symlink(goRoot, global.GvmGoroot)
	if err != nil {
		return err
	}

	output, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Print(string(output))

	return nil
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
