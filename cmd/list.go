/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/internal/vm"
)

var (
	IsRemote bool

	ListNumber int

	ListAll bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "gvm list go bin in local or remote",
	Long:  `gvm list go bin in local or remote`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	if !IsRemote {
		rd := vm.NewLocalVM()

		vs, err := rd.List()
		if err != nil {
			return err
		}

		goRoot, err := os.Readlink(global.GvmGoroot)
		if err != nil {
			// warning: do not use gvm activate before gvm list
			fmt.Printf("warning: do not use gvm activate before gvm list")
		}

		for _, v := range vs {
			if filepath.Join(global.GvmConfigDir, v) == goRoot {
				color.Blue("*\t%s\n", v)
			} else {
				fmt.Printf(" \t%s\n", v)
			}
		}
	} else {
		rd := vm.NewRemoteVM(&vm.RemoteVM{Registry: "https://go.dev/dl"})

		vs, err := rd.List()
		if err != nil {
			return err
		}
		if len(vs) == 0 {
			return errors.New("get remote version error")
		}

		if !ListAll {
			if len(vs) >= ListNumber {
				vs = vs[0:ListNumber]
			}
		}

		for _, v := range vs {
			fmt.Printf(" \t%s\n", v)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&IsRemote, "remote", "r", false, "is remote")
	listCmd.Flags().IntVarP(&ListNumber, "number", "n", 30, "show number")
	listCmd.Flags().BoolVarP(&ListAll, "all", "a", false, "is list all")
}
