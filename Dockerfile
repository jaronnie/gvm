FROM centos:7

MAINTAINER jaron@jaronnie.com

COPY ./dist/gvm_linux_amd64 /usr/bin/gvm

RUN yum -y install bash-completion \
    && gvm completion bash > /etc/bash_completion.d/gvm \
    && gvm init bash

