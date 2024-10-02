# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=ubuntu:jammy

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH=amd64
ARG LINUX_VERSION=6.1.45
ARG CONTEXT=linux

ARG GCC_VERSION=10
ENV GCC_VERSION=$GCC_VERSION
ARG USE_GCC_VERSION=$GCC_VERSION
ENV USE_GCC_VERSION=$USE_GCC_VERSION

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
        g++-${GCC_VERSION} \
        gawk \
        gcc-${GCC_VERSION} \
        gcc-${GCC_VERSION}-aarch64-linux-gnu \
        gcc-${GCC_VERSION}-arm-linux-gnueabi \
        gcc-${GCC_VERSION}-i686-linux-gnu \
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
        zstd \
    && \
    apt-get install -y --no-install-recommends \
        less \
        nano \
        openssh-server \
        screen \
        tmux \
        tree \
        vim \
    && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/aarch64-linux-gnu-gcc aarch64-linux-gnu-gcc /usr/bin/aarch64-linux-gnu-gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/arm-linux-gnueabi-gcc arm-linux-gnueabi-gcc /usr/bin/arm-linux-gnueabi-gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/i686-linux-gnu-gcc    i686-linux-gnu-gcc    /usr/bin/i686-linux-gnu-gcc-${GCC_VERSION} 100 && \
    gpg2 --locate-keys torvalds@kernel.org gregkh@kernel.org && \
    rm -rf /var/lib/apt/lists/*

# Prepare SSH for interactive debugging
RUN mkdir -p /run/sshd && \
    sed -E "s/#PermitRootLogin\b.*/PermitRootLogin yes/g" -i /etc/ssh/sshd_config && \
    sed -E "s/#PasswordAuthentication\b.*/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config

