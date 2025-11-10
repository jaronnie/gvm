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
	"runtime"
	"strings"

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

var SetUpGVMInWindows = `
# gvm shell setup
if (Test-Path "$env:USERPROFILE\gvm\.gvmrc.ps1") {
    . "$env:USERPROFILE\gvm\.gvmrc.ps1"
}
`

var GVMRCTemplateInWindows = `$env:GOROOT = "$env:USERPROFILE\gvm\goroot"
$env:PATH = "$env:PATH;$env:GOROOT\bin"
{{if empty (env "GOPATH") }}$env:GOPATH = "$env:USERPROFILE\gvm"
$env:GOBIN = "$env:GOPATH\bin"
$env:PATH = "$env:PATH;$env:GOBIN"{{else}}$env:GOPATH = "{{ env "GOPATH" }}"
$env:GOBIN = "$env:GOPATH\bin"
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
	if runtime.GOOS == "windows" {
		return initWindows(args)
	}
	return initUnix(args)
}

func initUnix(args []string) error {
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

func initWindows(args []string) error {
	shellType := "powershell"
	if len(args) == 1 {
		shellType = args[0]
	}

	fmt.Printf("🚀get shell type %s\n", shellType)

	var shellRcFile string
	var setupScript string
	var gvmrcFile string

	switch shellType {
	case "powershell", "pwsh":
		// PowerShell profile path
		output, err := exec.Command("powershell", "-Command", "echo $PROFILE").Output()
		if err == nil {
			shellRcFile = strings.TrimSpace(string(output))
		} else {
			// Fallback to default PowerShell profile location
			shellRcFile = filepath.Join(global.HomeDir, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
		}
		setupScript = SetUpGVMInWindows
		gvmrcFile = filepath.Join(global.GvmConfigDir, ".gvmrc.ps1")
	default:
		return errors.Errorf("unsupported shell type on Windows: %s. Please use 'powershell' or 'pwsh'", shellType)
	}

	// Ensure the profile directory exists
	profileDir := filepath.Dir(shellRcFile)
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return err
	}

	shellConfigFile, err := os.OpenFile(shellRcFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer shellConfigFile.Close()

	shellRcData, err := os.ReadFile(shellRcFile)
	if err != nil {
		return err
	}

	if !bytes.Contains(shellRcData, []byte("gvm shell setup")) {
		_, err = shellConfigFile.Write([]byte(setupScript))
		if err != nil {
			return err
		}
	}

	// go env GOPATH
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		output, _ := exec.Command("go", "env", "GOPATH").Output()
		output = bytes.TrimSpace(output)
		_ = os.Setenv("GOPATH", string(output))
	}

	template, err := utilx.ParseTemplate(nil, []byte(GVMRCTemplateInWindows))
	if err != nil {
		return err
	}

	err = os.WriteFile(gvmrcFile, template, 0o644)
	if err != nil {
		return err
	}

	// cp gvm exec binary to %USERPROFILE%/gvm/bin
	path, _ := utilx.LookPath(rootCmd.Use)
	if path == "" {
		// 如果找不到 gvm, 则复制当前二进制文件到 %USERPROFILE%/gvm/bin
		fileStat, err := os.Stat(os.Args[0])
		if err != nil {
			return err
		}
		file, err := os.ReadFile(os.Args[0])
		if err != nil {
			return err
		}
		_ = os.MkdirAll(filepath.Join(global.GvmConfigDir, "bin"), 0o755)
		binaryName := "gvm.exe"
		err = os.WriteFile(filepath.Join(global.GvmConfigDir, "bin", binaryName), file, fileStat.Mode())
		if err != nil {
			return err
		}
	}

	fmt.Printf("🚀please restart your PowerShell or run `. $PROFILE` to activate gvm\n")

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
