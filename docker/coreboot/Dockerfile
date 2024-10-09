# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG TARGETARCH=amd64

#=============
# "base" stage with all needed build dependencies
ARG SOURCE_IMAGE=ubuntu:jammy
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH

ARG COREBOOT_VERSION=4.19
ENV COREBOOT_VERSION=${COREBOOT_VERSION}
ARG CONTEXT=coreboot

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_COREBOOT_VERSION=${COREBOOT_VERSION}

ENV DEBIAN_FRONTEND=noninteractive
# Directory for coreboot toolchain, MEAnalyser etc.
ENV TOOLSDIR=/tools
# Use coreboot mirror
ENV BUILDGCC_OPTIONS=-m

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        acpica-tools \
        bc \
        bison \
        bsdmainutils \
        build-essential \
        curl \
        flex \
        git \
        gnat \
        imagemagick \
        libelf-dev \
        libncurses5-dev \
        libnss3-dev \
        libssl-dev \
        m4 \
        nasm \
        nodejs \
        openssh-client \
        pkgconf \
        python-is-python3 \
        python3-pip \
        qemu-system-x86 \
        tzdata \
        uuid-dev \
        zlib1g-dev \
    && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        apt-get install -y --no-install-recommends \
            iucode-tool; \
    fi; \
    apt-get install -y --no-install-recommends \
        less \
        nano \
        openssh-server \
        screen \
        tmux \
        tree \
        vim \
    && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}"

# Let coreboot know the toolchain path
ENV XGCCPATH=${TOOLSDIR}/coreboot/util/crossgcc/xgcc/bin/

# Add pre-compiled coreboot toolchain and utils
COPY coreboot-${COREBOOT_VERSION}/xgcc-${TARGETARCH} ${TOOLSDIR}/coreboot/util/crossgcc/xgcc
COPY coreboot-${COREBOOT_VERSION}/utils-${TARGETARCH}/* /usr/local/bin/
RUN git clone --depth 1 https://github.com/platomav/MEAnalyzer.git ${TOOLSDIR}/MEAnalyzer/

# Prepare SSH for interactive debugging
RUN mkdir -p /run/sshd && \
    sed -E "s/#PermitRootLogin\b.*/PermitRootLogin yes/g" -i /etc/ssh/sshd_config && \
    sed -E "s/#PasswordAuthentication\b.*/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config

