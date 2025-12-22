# gvm

[![GitHub release](https://img.shields.io/github/release/jaronnie/gvm.svg?style=flat-square)](https://github.com/jaronnie/gvm/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jaronnie/gvm/ci.yaml?branch=main&label=gvm-golint&logo=github&style=flat-square)](https://github.com/jaronnie/gvm/actions?query=workflow%3Agvm-golint)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jaronnie/gvm/ci.yaml?branch=main&label=goreleaser-gvm&logo=github&style=flat-square)](https://github.com/jaronnie/gvm/actions?query=workflow%3Agoreleaser-gvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaronnie/gvm?style=flat-square)](https://goreportcard.com/report/github.com/jaronnie/gvm)
[![codecov](https://img.shields.io/codecov/c/github/jaronnie/gvm?logo=codecov&style=flat-square)](https://codecov.io/gh/jaronnie/gvm)

golang 多版本管理工具

## 下载 gvm

### Docker

```shell
docker run -it ghcr.io/jaronnie/gvm:latest bash
```

### 从源码下载(需要依赖 go 环境)

```shell
go install github.com/jaronnie/gvm@latest
```

### 下载二进制文件

[点击下载](https://github.com/jaronnie/gvm/releases)

#### Linux/macOS

```shell
# linux amd64
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.9.0/gvm_1.8.0_Linux_x86_64.tar.gz
# darwin amd64
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.9.0/gvm_1.8.0_Darwin_x86_64.tar.gz
# darwin arm64(m 系列)
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.9.0/gvm_1.8.0_Darwin_arm64.tar.gz
```

```shell
tar -zxvf gvm.tar.gz
mv gvm /usr/local/bin
```

#### Windows

```powershell
# 下载 Windows 版本
# amd64
Invoke-WebRequest -Uri "https://github.com/jaronnie/gvm/releases/download/v1.9.0/gvm_1.8.0_Windows_x86_64.zip" -OutFile "gvm.zip"
# arm64
Invoke-WebRequest -Uri "https://github.com/jaronnie/gvm/releases/download/v1.9.0/gvm_1.8.0_Windows_arm64.zip" -OutFile "gvm.zip"

# 解压
Expand-Archive -Path gvm.zip -DestinationPath $env:USERPROFILE\gvm\bin
# 添加到 PATH (需要管理员权限，或者手动添加到系统环境变量)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\gvm\bin", "User")
```

## 使用 gvm

### 第一步: 初始化 gvm

#### Linux/macOS

```shell
gvm init
```

**让环境变量生效**

```shell
# zsh
source ~/.zshrc
# bash
source ~/.bashrc
```

#### Windows

```powershell
# 初始化 gvm (使用 PowerShell)
gvm init powershell

# 让环境变量生效 (重启 PowerShell 或运行)
. $PROFILE
```

**注意**: Windows 上使用 gvm activate 命令可能需要管理员权限来创建符号链接。如果没有管理员权限，建议启用 Windows 10/11 的开发者模式。

### 第二步: gvm 命令补全(可选)

#### Linux/macOS

```shell
# zsh:
# echo "autoload -U compinit; compinit" >> ~/.zshrc
gvm completion zsh > "${fpath[1]}/_gvm"

# linux bash:
gvm completion bash > /etc/bash_completion.d/gvm
```

#### Windows

```powershell
# PowerShell 补全支持
gvm completion powershell > $env:USERPROFILE\Documents\WindowsPowerShell\Modules\gvm\gvm.ps1
```

### 第三步: 下载 go 版本

```shell
# will install go 1.23.5
gvm install go1.23.5

# install offline, default package file path is $HOME/gvm
gvm install go1.23.5 --offline

# install offline, package file path is .
gvm install go1.23.5 --offline -p .
```

### 第四步: 列举下载的 go 版本

```shell
# list local go
gvm list

# list remote go
gvm list --remote

# list remote go with limit number
gvm list --remote -n 100

# list all
gvm list --remote --all
```

### 第五步: 激活 go 版本

```shell
# will activate go 1.23.5 environment
gvm activate go1.23.5
```

### 第六步: 卸载 go 版本

```shell
gvm uninstall go1.23.5
```

### go 版本下载路径

#### Linux/macOS
```shell
ls $HOME/gvm/go*
```

#### Windows
```powershell
ls $env:USERPROFILE\gvm\go*
```

## 问题与解决

### windows: 无法加载文件 C:\Users\xx\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1，因为在此系统上禁止运行脚本

```shell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUse
```

### 下载 go 版本失败, 设置代理

```shell
gvm config --registry https://mirrors.aliyun.com/golang
```

### bash completion error: bash: _get_comp_words_by_ref: command not found

```shell
# centos
yum -y install bash-completion
```

### bash: permission denied: /etc/bash_completion.d/gvm

```shell
gvm completion bash | sudo tee /etc/bash_completion.d/gvm > /dev/null
```
