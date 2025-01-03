# gvm

[![GitHub release](https://img.shields.io/github/release/jaronnie/gvm.svg?style=flat-square)](https://github.com/jaronnie/gvm/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jaronnie/gvm/ci.yaml?branch=main&label=gvm-golint&logo=github&style=flat-square)](https://github.com/jaronnie/gvm/actions?query=workflow%3Agvm-golint)
[![Build Status](https://img.shields.io/github/actions/workflow/status/jaronnie/gvm/ci.yaml?branch=main&label=goreleaser-gvm&logo=github&style=flat-square)](https://github.com/jaronnie/gvm/actions?query=workflow%3Agoreleaser-gvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaronnie/gvm?style=flat-square)](https://goreportcard.com/report/github.com/jaronnie/gvm)
[![codecov](https://img.shields.io/codecov/c/github/jaronnie/gvm?logo=codecov&style=flat-square)](https://codecov.io/gh/jaronnie/gvm)

golang version manage

## quick start

```shell
docker run -it jaronnie/gvm:latest bash
# docker run -it jaronnie/gvm:latest-arm64 bash
gvm install go1.20
gvm activate go1.20
```

## download gvm

### download by source

```shell
go install github.com/jaronnie/gvm@latest
```

### download from releases

[download](https://github.com/jaronnie/gvm/releases)

```shell
curl -L -o gvm.tar.gz https://github.com/jaronnie/gvm/releases/download/v1.6.0/gvm_1.6.0_Linux_x86_64.tar.gz
```

## gvm init

```shell
gvm init
# gvm init <shellType>
```

## gvm complete

```shell
zsh:
# echo "autoload -U compinit; compinit" >> ~/.zshrc
gvm completion zsh > "${fpath[1]}/_gvm"

linux bash:
gvm completion bash > /etc/bash_completion.d/gvm
```

## gvm install

```shell
# will install go 1.18
gvm install go1.18

# will install go 1.18.5 version
gvm install go1.18.5

# install offline, default package file path is ~/gvm
gvm install go1.18.5 --offline

# install offline, package file path is .
gvm install go1.18.5 --offline -p .
```

## gvm list

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

## gvm activate

```shell
# will activate go 1.18 environment
gvm activate go1.18
```

## gvm uninstall

```shell
gvm uninstall go1.18.5
```

## QA

### download go package error

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

## todo

- [ ] support windows
