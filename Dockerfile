FROM centos:7.9.2009

LABEL MAINTAINER jaron@jaronnie.com

COPY ./dist/gvm_linux_amd64/gvm /usr/bin/gvm

RUN yum -y install bash-completion \
    && gvm completion bash > /etc/bash_completion.d/gvm \
    && gvm init bash

