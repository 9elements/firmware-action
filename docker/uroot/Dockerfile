# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG GO_VERSION=1.22
ARG SOURCE_IMAGE=golang:${GO_VERSION}

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG UROOT_VERSION=0.14.0
ARG CONTEXT=uroot

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_UROOT_VERSION=v${UROOT_VERSION}

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        less \
        nano \
        openssh-server \
        screen \
        tmux \
        tree \
        vim \
    && \
    rm -rf /var/lib/apt/lists/*

# Prepare SSH for interactive debugging
RUN mkdir -p /run/sshd && \
    sed -E "s/#PermitRootLogin\b.*/PermitRootLogin yes/g" -i /etc/ssh/sshd_config && \
    sed -E "s/#PasswordAuthentication\b.*/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config

