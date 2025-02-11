# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=ubuntu:noble
ARG TARGETARCH=amd64

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH
ARG UBOOT_VERSION=2025.01
ARG CONTEXT=uboot

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_UBOOT_VERSION=v${UBOOT_VERSION}

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bc \
        bison \
        build-essential \
        clang \
        coccinelle \
        device-tree-compiler \
        dfu-util \
        efitools \
        flex \
        gcc \
        gdisk \
        git \
        graphviz \
        imagemagick \
        libgnutls28-dev \
        libguestfs-tools \
        liblz4-tool \
        libncurses-dev \
        libpython3-dev \
        libsdl2-dev \
        libssl-dev \
        lz4 \
        lzma \
        lzma-alone \
        openssl \
        pkg-config \
        python3 \
        python3-asteval \
        python3-coverage \
        python3-filelock \
        python3-pkg-resources \
        python3-pycryptodome \
        python3-pyelftools \
        python3-pytest \
        python3-pytest-xdist \
        python3-sphinx-rtd-theme \
        python3-sphinxcontrib.apidoc \
        python3-subunit \
        python3-testtools \
        python3-virtualenv \
        swig \
        uuid-dev \
    && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        apt-get install -y --no-install-recommends \
            gcc-aarch64-linux-gnu; \
    fi && \
    apt-get install -y --no-install-recommends \
        less \
        nano \
        screen \
        tmux \
        tree \
        vim \
    && \
    rm -rf /var/lib/apt/lists/*

