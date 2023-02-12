/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/internal/vm"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "gvm list go bin in local or remote",
	Long:  `gvm list go bin in local or remote`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	rd := vm.NewReadDirVM()

	vs, err := rd.List()
	if err != nil {
		return err
	}

	goRoot, err := os.Readlink(global.GVM_GOROOT)
	if err != nil {
		// warning: do not use gvm activate before gvm list
	}

	for _, v := range vs {
		if filepath.Join(global.GVM_CONFIG_DIR, v) == goRoot {
			color.Blue("*\t%s\n", v)
		} else {
			fmt.Printf(" \t%s\n", v)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
