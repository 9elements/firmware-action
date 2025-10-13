# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=ubuntu:noble
ARG TARGETARCH=amd64

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH
ARG LINUX_VERSION=6.1.127
ARG LINUX_VERSION_MAJOR_MINOR=6.1
ARG CONTEXT=linux

ARG GCC_VERSION=10
ENV GCC_VERSION=$GCC_VERSION
ARG USE_GCC_VERSION=$GCC_VERSION
ENV USE_GCC_VERSION=$USE_GCC_VERSION

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_LINUX_VERSION=$LINUX_VERSION_MAJOR_MINOR

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
        g++-${GCC_VERSION} \
        gawk \
        gcc-${GCC_VERSION} \
        gcc-${GCC_VERSION}-arm-linux-gnueabi \
        gcc-${GCC_VERSION}-i686-linux-gnu \
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
        zstd \
    && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        apt-get install -y --no-install-recommends \
            gcc-${GCC_VERSION}-aarch64-linux-gnu; \
    fi && \
    if [ "${TARGETARCH}" = 'arm64' ]; then \
        apt-get install -y --no-install-recommends \
            gcc-${GCC_VERSION}-x86-64-linux-gnu; \
    fi && \
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
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/aarch64-linux-gnu-gcc aarch64-linux-gnu-gcc /usr/bin/aarch64-linux-gnu-gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/arm-linux-gnueabi-gcc arm-linux-gnueabi-gcc /usr/bin/arm-linux-gnueabi-gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/i686-linux-gnu-gcc    i686-linux-gnu-gcc    /usr/bin/i686-linux-gnu-gcc-${GCC_VERSION} 100; \
    gpg2 --locate-keys torvalds@kernel.org gregkh@kernel.org && \
    rm -rf /var/lib/apt/lists/*

