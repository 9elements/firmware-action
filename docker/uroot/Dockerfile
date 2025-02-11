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
        curl \
        git \
        less \
        nano \
        screen \
        tmux \
        tree \
        vim \
        wget \
    && \
    rm -rf /var/lib/apt/lists/*

