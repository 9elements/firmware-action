# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

#=============
# "base" stage with all needed build dependencies
FROM ubuntu:jammy AS base

ARG TARGETARCH=amd64
ARG COREBOOT_VERSION=4.19
ARG CONTEXT=coreboot

# Verification test
ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_COREBOOT_VERSION=$COREBOOT_VERSION

ENV DEBIAN_FRONTEND=noninteractive
# Directory for coreboot toolchain, MEAnalyser etc.
ENV TOOLSDIR=/tools
# Use coreboot mirror
ENV BUILDGCC_OPTIONS=-m

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bc \
        bison \
        build-essential \
        curl \
        flex \
        git \
        gnat \
        iucode-tool \
        libelf-dev \
        libncurses5-dev \
        libssl-dev \
        m4 \
        nodejs \
        pkgconf \
        python-is-python3 \
        python3-pip \
        tzdata \
        zlib1g-dev \
    && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    pip3 install --no-cache-dir \
        colorama==0.4.6 \
        crccheck==1.3.0 \
        pltable==1.0.2 \
    && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}"


#=============
# "toolchain" stage to build the coreboot toolchain
FROM base AS toolchain

# Compile coreboot toolchain
WORKDIR $TOOLSDIR
RUN git clone "https://review.coreboot.org/coreboot.git"
WORKDIR $TOOLSDIR/coreboot
RUN git checkout "${COREBOOT_VERSION}"
RUN make crossgcc-aarch64 CPUS="$(nproc)" && \
    make crossgcc-arm CPUS="$(nproc)" && \
    make crossgcc-i386 CPUS="$(nproc)" && \
    make crossgcc-x64 CPUS="$(nproc)"
WORKDIR $TOOLSDIR/coreboot
RUN make -C util/ifdtool install && \
    make -C util/cbfstool install

# Install MEAnalyser
WORKDIR $TOOLSDIR
RUN git clone https://github.com/platomav/MEAnalyzer.git


#=============
# "final" stage is the actual product with everything included
FROM base AS final

# Let coreboot know the toolchain path
ENV XGCCPATH=$TOOLSDIR/coreboot/util/crossgcc/xgcc/bin/
RUN echo "${XGCCPATH}"

# Copy over things from previous stage(s)
COPY --from=toolchain $XGCCPATH/.. $XGCCPATH/..
COPY --from=toolchain $TOOLSDIR/MEAnalyzer $TOOLSDIR/

