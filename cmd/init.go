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
	// Detect and show existing Go installations
	existingInstallations, _ := detectExistingGoInstallations()
	if existingInstallations != "" {
		fmt.Printf("❌  Initialization aborted due to existing Go installation conflicts %s.\n", existingInstallations)
		return nil
	}

	if runtime.GOOS == "windows" {
		return initWindows(args)
	}
	return initUnix(args)
}

func initUnix(args []string) error {
	shellType := os.Getenv("SHELL")

	// If shell type is provided as argument, use it
	if len(args) == 1 {
		shellType = args[0]
	}

	// If still no shell type, try to detect from common locations
	if shellType == "" {
		// Try to detect shell from common shell paths
		if _, err := os.Stat("/bin/bash"); err == nil {
			shellType = "bash"
		} else if _, err = os.Stat("/bin/sh"); err == nil {
			shellType = "sh"
		} else if _, err = os.Stat("/bin/zsh"); err == nil {
			shellType = "zsh"
		} else if _, err = os.Stat("/bin/fish"); err == nil {
			shellType = "fish"
		} else {
			return errors.New("cannot determine shell type. Please specify shell type: gvm init <bash|zsh|sh>")
		}
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
	case "fish":
		shellRcFile = filepath.Join(global.HomeDir, ".fishrc")
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

	fmt.Print("🚀gvm initialized successfully!\n")
	fmt.Printf("Please run `source %s` to activate gvm\n", shellRcFile)

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

	fmt.Print("🚀GVM initialized successfully!\n")
	fmt.Printf("Please restart your PowerShell or run `. %s` to activate GVM\n", shellRcFile)

	return nil
}

// detectExistingGoInstallations detects existing Go installations not managed by GVM
func detectExistingGoInstallations() (string, error) {
	// Check PATH for go executable
	_, err := exec.LookPath("go")
	if err != nil {
		// Go not found in PATH, no conflicts
		return "", nil
	}

	// Get the actual GOROOT from the existing Go installation
	output, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
	if err != nil {
		return "", nil // If we can't get GOROOT, assume no conflict
	}

	existingGoroot := strings.TrimSpace(string(output))
	if existingGoroot == "" {
		return "", nil
	}

	// Check if this Go installation is managed by GVM
	// GVM-managed installations are in ~/gvm/go{version} format
	if !strings.HasPrefix(existingGoroot, global.GvmConfigDir) {
		return existingGoroot, nil
	}
	return "", nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
