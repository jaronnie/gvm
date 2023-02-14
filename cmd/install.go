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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/jaronnie/gvm/internal/global"
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
	RunE:  install,
}

func install(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	gov := args[0]
	if !strings.HasPrefix(gov, "go") {
		return errors.New("invalid go version, please use gox.x")
	}

	fmt.Printf("🖕Install go %s\n", gov)

	var packagePath string

	if !IsOffline {
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

		if resp.ContentLength <= 0 {
			return errors.New("content length less than 0")
		}

		filename := filepath.Base(installUrl)
		file, err := os.Create(filepath.Join(global.GVM_CONFIG_DIR, filename))
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

		bar.Finish()

		fmt.Printf("🔥Install %s successfully\n", gov)
	} else {
		packagePath = filepath.Join(PackagePath, fmt.Sprintf("%s.%s-%s.tar.gz", gov, runtime.GOOS, runtime.GOARCH))
	}

	fmt.Printf("🚀Start to untar %s to %s\n", packagePath, global.GVM_CONFIG_DIR)

	err := utilx.Untargz(packagePath, global.GVM_CONFIG_DIR)
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

	installCmd.Flags().BoolVarP(&IsOffline, "offline", "", false, "is offline")

	homeDir, _ := os.UserHomeDir()
	installCmd.Flags().StringVarP(&PackagePath, "package-path", "p", filepath.Join(homeDir, GVM), "if use offline flag, set package path")
}
