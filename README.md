# gvm

golang version manage

## 配置

```shell
# ~/.zshrc
export GOROOT=$HOME/.gvm/goroot
export PATH=$PATH:$GOROOT/bin
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

## install

```shell
# will install go 1.18
gvm install go1.18

# will install go 1.18.5 version
gvm install go1.18.5
```

## activate

```shell
# will activate go 1.18 environment
gvm activate go1.18
```

## todo

- [ ] gvm init