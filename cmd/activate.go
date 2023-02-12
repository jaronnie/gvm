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
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
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

	goRoot := filepath.Join(global.GVM_CONFIG_DIR, gov)

	_ = os.Remove(global.GVM_GOROOT)

	err := os.Symlink(goRoot, global.GVM_GOROOT)
	if err != nil {
		return err
	}

	output, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf(string(output))

	return nil
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
