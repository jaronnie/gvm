/*
Copyright © 2023 jaronnie jaron@jaronnie.com

*/

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/jaronnie/gvm/internal/global"
	"github.com/jaronnie/gvm/internal/vm"
	"github.com/jaronnie/gvm/utilx"
)

var (
	IsOffline   bool
	PackagePath string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install go bin",
	Long:  `install go bin`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		rd := vm.NewRemoteVM(&vm.RemoteVM{
			Registry: "https://go.dev/dl",
			Cache:    Cache,
		})

		vs, err := rd.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return vs, cobra.ShellCompDirectiveDefault
	},
	RunE: install,
}

func install(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	rd := vm.NewLocalVM()

	vs, _ := rd.List()
	if lo.Contains(vs, gov) {
		return errors.Errorf("go %s has been installed", gov)
	}
	fmt.Printf("🖕Install go %s\n", gov)

	var packagePath string

	var archiveExt string
	if runtime.GOOS == "windows" {
		archiveExt = ".zip"
	} else {
		archiveExt = ".tar.gz"
	}

	if !IsOffline {
		registry := viper.GetString("registry")
		if registry == "" {
			registry = Registry
		}

		installURL := fmt.Sprintf("%s/%s.%s-%s%s", registry, gov, runtime.GOOS, runtime.GOARCH, archiveExt)
		fmt.Printf("🌿Install from %s\n", installURL)

		resp, err := utilx.GetRawResponse(installURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.ContentLength <= 0 {
			return errors.New("content length less than 0")
		}

		filename := filepath.Base(installURL)
		file, err := os.Create(filepath.Join(global.GvmConfigDir, filename))
		if err != nil {
			return err
		}
		defer file.Close()

		packagePath = file.Name()

		defer os.Remove(file.Name())

		length := resp.Header.Get("Content-Length")
		lengthKb, _ := strconv.ParseFloat(length, 64)

		fmt.Printf("🚀Save to: %s\n", file.Name())

		bar := pb.New(int(lengthKb)).SetUnits(pb.U_BYTES_DEC).SetRefreshRate(time.Millisecond * 10)

		// 显示下载速度
		bar.ShowSpeed = true
		// 显示剩余时间
		bar.ShowTimeLeft = true
		// 显示完成时间
		bar.ShowFinalTime = true
		bar.SetWidth(80)
		bar.Start()

		writer := io.MultiWriter(file, bar)

		_, err = io.Copy(writer, resp.Body)
		if err != nil {
			return err
		}

		bar.Finish()

		fmt.Printf("🔥Install %s successfully\n", gov)
	} else {
		packagePath = filepath.Join(PackagePath, fmt.Sprintf("%s.%s-%s%s", gov, runtime.GOOS, runtime.GOARCH, archiveExt))
	}

	fmt.Printf("🚀Start to extract %s to %s\n", packagePath, global.GvmConfigDir)

	var err error
	if runtime.GOOS == "windows" {
		err = utilx.Unzip(packagePath, global.GvmConfigDir)
	} else {
		err = utilx.Untargz(packagePath, global.GvmConfigDir)
	}
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(global.GvmConfigDir, "go"), filepath.Join(global.GvmConfigDir, gov))
	if err != nil {
		return err
	}

	// chmod -R 555 gox.xx
	err = os.Chmod(filepath.Join(global.GvmConfigDir, gov), 0o544)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&IsOffline, "offline", "", false, "is offline")

	homeDir, _ := os.UserHomeDir()
	installCmd.Flags().StringVarP(&PackagePath, "package-path", "p", filepath.Join(homeDir, GVM), "if use offline flag, set package path")
}
