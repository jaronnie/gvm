/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/utilx"
	"github.com/jaronnie/gvm/utilx/downloader"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install go bin",
	Long:  `install go bin`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:  install,
}

var p *tea.Program

func install(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	fmt.Printf("🖕Install go %s\n", gov)

	registry := viper.GetString("registry")
	if registry == "" {
		registry = Registry
	}

	installUrl := fmt.Sprintf("%s/%s.%s-%s.tar.gz", registry, gov, runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🌿Install from %s\n", installUrl)

	resp, err := utilx.GetRawResponse(installUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Don't add TUI if the header doesn't include content size
	// it's impossible see progress without total
	if resp.ContentLength <= 0 {
		return errors.New("content length less than 0")
	}

	filename := filepath.Base(installUrl)
	file, err := os.Create(filepath.Join(global.GVM_CONFIG_DIR, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	defer os.Remove(file.Name())

	pw := &downloader.ProgressWriter{
		Total:  int(resp.ContentLength),
		File:   file,
		Reader: resp.Body,
		OnProgress: func(ratio float64) {
			p.Send(downloader.ProgressMsg(ratio))
		},
	}

	m := downloader.Model{
		Pw:       pw,
		Progress: progress.New(progress.WithDefaultGradient()),
	}
	// Start Bubble Tea
	p = tea.NewProgram(m)

	// Start the download
	go pw.Start()

	if _, err := p.Run(); err != nil {
		return err
	}

	// TODO
	// 断点下载

	fmt.Printf("🔥Install %s successfully\n", gov)
	fmt.Printf("🚀Start to untar %s to %s\n", file.Name(), global.GVM_CONFIG_DIR)

	err = utilx.Untargz(file.Name(), global.GVM_CONFIG_DIR)
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(global.GVM_CONFIG_DIR, "go"), filepath.Join(global.GVM_CONFIG_DIR, gov))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
