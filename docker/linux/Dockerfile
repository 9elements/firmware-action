# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=ubuntu:jammy

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH=amd64
ARG LINUX_VERSION=6.1.45
ARG CONTEXT=linux

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_LINUX_VERSION=$LINUX_VERSION

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        autoconf \
        bc \
        bison \
        build-essential \
        ca-certificates \
        dirmngr \
        dkms \
        flex \
        gawk \
        git \
        gnupg2 \
        libelf-dev \
        libiberty-dev \
        libncurses-dev \
        libpci-dev \
        libssl-dev \
        libudev-dev \
        llvm \
        nodejs \
        openssl \
        wget \
    && \
    gpg2 --locate-keys torvalds@kernel.org gregkh@kernel.org && \
    rm -rf /var/lib/apt/lists/*
