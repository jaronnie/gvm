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

```shell
# linux amd64
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.8.0/gvm_1.8.0_Linux_x86_64.tar.gz
# darwin amd64
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.8.0/gvm_1.8.0_Darwin_x86_64.tar.gz
# darwin arm64(m 系列)
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.8.0/gvm_1.8.0_Darwin_arm64.tar.gz
```

```shell
tar -zxvf gvm.tar.gz
mv gvm /usr/local/bin
```

## 使用 gvm

### 第一步: 初始化 gvm

```shell
gvm init
# 如果执行失败, 手动指定 shell 类型
# gvm init bash
```

**让环境变量生效**

```shell
# zsh
source ~/.zshrc
# bash
source ~/.bashrc
```

### 第二步: gvm 命令补全(可选)

```shell
zsh:
# echo "autoload -U compinit; compinit" >> ~/.zshrc
gvm completion zsh > "${fpath[1]}/_gvm"

linux bash:
gvm completion bash > /etc/bash_completion.d/gvm
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

## gvm 相关配置

### 环境变量设置

```shell
cat $HOME/gvm/.gvmrc

# get env
export GOROOT=$HOME/gvm/goroot
export PATH=$PATH:$GOROOT/bin
export GOPATH=$HOME/gvm
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

### go 版本下载路径

```shell
ls $HOME/gvm/go*
```

## 问题与解决

### 下载 go 版本失败, 设置代理

```shell
gvm config --registry https://mirrors.aliyun.com/golang
```

### bash completion error: bash: _get_comp_words_by_ref: command not found

```shell
# centos
yum -y install bash-completion
```

### gvm init: Error: can not get shell type

```shell
# gvm init <shellType>
gvm init bash
```

### bash: permission denied: /etc/bash_completion.d/gvm

```shell
gvm completion bash | sudo tee /etc/bash_completion.d/gvm > /dev/null
```
