/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/utilx"
)

var (
	SetUpGVMInUnix = `
# gvm shell setup
if [ -f "${HOME}/gvm/.gvmrc" ]; then
    source "${HOME}/gvm/.gvmrc"
fi
`
)

var (
	GVMRCTemplateInUnix = `export GOROOT=$HOME/gvm/goroot
export PATH=$PATH:$GOROOT/bin
{{if empty (env "GOPATH") }}export GOPATH=$HOME/gvm
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN{{else}}export GOPATH={{ env "GOPATH" }}
export GOBIN=$GOPATH/bin
{{end}}
`
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "gvm init",
	Long:  `gvm init`,
	RunE:  initx,
}

func initx(cmd *cobra.Command, args []string) error {
	shellType := os.Getenv("SHELL")

	if len(args) == 1 {
		shellType = args[0]
	}

	if shellType == "" {
		return errors.New("can not get shell type")
	}

	fmt.Printf("🚀get shell type %s\n", shellType)

	// get shell rc file
	var shellRcFile string
	switch filepath.Base(shellType) {
	case "sh":
		shellRcFile = filepath.Join(global.HOME_DIR, ".shrc")
	case "bash":
		shellRcFile = filepath.Join(global.HOME_DIR, ".bashrc")
	case "zsh":
		shellRcFile = filepath.Join(global.HOME_DIR, ".zshrc")
	}

	shellConfigfile, err := os.OpenFile(shellRcFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		return err
	}
	defer shellConfigfile.Close()

	shellRcData, err := os.ReadFile(shellRcFile)
	if err != nil {
		return err
	}

	if !bytes.Contains(shellRcData, []byte("gvm shell setup")) {
		_, err = shellConfigfile.Write([]byte(SetUpGVMInUnix))
		if err != nil {
			return err
		}
	}

	// go env GOPATH
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		output, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			// warning
			output = []byte("")
		} else {
			output = bytes.TrimRight(output, "\n")
			_ = os.Setenv("GOPATH", string(output))
		}
	}

	template, err := utilx.ParseTemplate(nil, []byte(GVMRCTemplateInUnix))
	if err != nil {
		return err
	}

	err = os.WriteFile(global.GVM_CONFIG_RC, template, 0744)
	if err != nil {
		return err
	}

	fmt.Printf("🚀please exec `source %s` to activate gvm\n", shellRcFile)

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
