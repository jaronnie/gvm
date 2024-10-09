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

var SetUpGVMInUnix = `
# gvm shell setup
if [ -f "${HOME}/gvm/.gvmrc" ]; then
    source "${HOME}/gvm/.gvmrc"
fi
`

var GVMRCTemplateInUnix = `export GOROOT=$HOME/gvm/goroot
export PATH=$PATH:$GOROOT/bin
{{if empty (env "GOPATH") }}export GOPATH=$HOME/gvm
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN{{else}}export GOPATH={{ env "GOPATH" }}
export GOBIN=$GOPATH/bin
{{end}}
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <shellType>",
	Short: "gvm init",
	Long:  `gvm init`,
	RunE:  initx,
}

func initx(cmd *cobra.Command, args []string) error {
	shellType := os.Getenv("SHELL")
	fmt.Printf("get SHELL env: %s\n", shellType)

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
		shellRcFile = filepath.Join(global.HomeDir, ".shrc")
	case "bash":
		shellRcFile = filepath.Join(global.HomeDir, ".bashrc")
	case "zsh":
		shellRcFile = filepath.Join(global.HomeDir, ".zshrc")
	}

	shellConfigFile, err := os.OpenFile(shellRcFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o744)
	if err != nil {
		return err
	}
	defer shellConfigFile.Close()

	shellRcData, err := os.ReadFile(shellRcFile)
	if err != nil {
		return err
	}

	if !bytes.Contains(shellRcData, []byte("gvm shell setup")) {
		_, err = shellConfigFile.Write([]byte(SetUpGVMInUnix))
		if err != nil {
			return err
		}
	}

	// go env GOPATH
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		output, _ := exec.Command("go", "env", "GOPATH").Output()
		output = bytes.TrimRight(output, "\n")
		_ = os.Setenv("GOPATH", string(output))
	}

	template, err := utilx.ParseTemplate(nil, []byte(GVMRCTemplateInUnix))
	if err != nil {
		return err
	}

	err = os.WriteFile(global.GvmConfigRc, template, 0o744)
	if err != nil {
		return err
	}

	// cp gvm exec binary to $HOME/gvm/bin
	path, _ := utilx.LookPath(rootCmd.Use)
	if path == "" {
		// 如果找不到 gvm, 则复制当前二进制文件到 $HOME/gvm/bin
		fileStat, err := os.Stat(os.Args[0])
		if err != nil {
			return err
		}
		file, err := os.ReadFile(os.Args[0])
		if err != nil {
			return err
		}
		_ = os.MkdirAll(filepath.Join(global.GvmConfigDir, "bin"), 0o755)
		err = os.WriteFile(filepath.Join(global.GvmConfigDir, "bin", "gvm"), file, fileStat.Mode())
		if err != nil {
			return err
		}
	}

	fmt.Printf("🚀please exec `source %s` to activate gvm\n", shellRcFile)

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
