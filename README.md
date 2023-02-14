# gvm

golang version manage

`one issue need help: can not use in goland's terminal!`

## init

初始化 gvm
```shell
gvm init
# gvm init <shellType>
```

## auto completion

```shell
zsh:
# echo "autoload -U compinit; compinit" >> ~/.zshrc
gvm completion zsh > "${fpath[1]}/_gvm"

linux bash:
gvm completion bash > /etc/bash_completion.d/gvm
```

## install go

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

## list

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

## activate go

```shell
# will activate go 1.18 environment
gvm activate go1.18
```

## uninstall go

```shell
gvm uninstall go1.18.5
```

## QA

### 下载 go 压缩包失败

可以通过设置阿里云的 registry 镜像源解决
```shell
gvm config --registry https://mirrors.aliyun.com/golang
```

### bash 补全错误 bash: _get_comp_words_by_ref: command not found

```shell
# centos
yum -y install bash-completion
```

### gvm init: Error: can not get shell type

```shell
# gvm init <shellType>
gvm init bash
```


## todo

- [ ] support windows