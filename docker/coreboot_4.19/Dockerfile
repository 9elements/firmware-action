FROM ubuntu:22.04

ARG TARGETARCH=amd64

# Install dependencies
ENV DEBIAN_FRONTEND=noninteractive
ENV BUILDGCC_OPTIONS=-m

# Install git from PPA as v2.34 cannot clone coreboot...
RUN apt-get -y update && \
    apt-get -y upgrade && \
    apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get install -y git curl tzdata \
    bison build-essential curl flex git gnat iucode-tool \
    libncurses5-dev m4 zlib1g-dev bc libelf-dev python3-pip python-is-python3 && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    pip3 install --no-cache-dir colorama crccheck pltable

# coreboot toolchain
RUN mkdir /repo && cd /repo && git clone https://review.coreboot.org/coreboot.git && \
    cd /repo/coreboot && \
    git checkout 4d661eebe95b06acae7d99777ed36ba56a560112 && \
    cd /repo/coreboot/util/crossgcc; make build-i386 CPUS=$(nproc) && make build_iasl && \
    make build_nasm && \
    cd /repo/coreboot && rm -rf -v !\("util/"\) && \
    cd /repo/coreboot/ && \
    make -C /repo/coreboot/util/ifdtool install && \
    make -C /repo/coreboot/util/cbfstool install && \
    rm -rf /repo/coreboot/util/crossgcc/tarballs/ && \
    rm -rf /repo/coreboot/.git && \
    rm -rf /repo/coreboot/src && rm -rf /repo/coreboot/payloads && rm -rf /repo/coreboot/tests

# Let coreboot know the toolchain path
ENV XGCCPATH=/repo/coreboot/util/crossgcc/xgcc/bin/

# Install MEAnalyser
RUN mkdir -p /repo && cd repo && git clone https://github.com/platomav/MEAnalyzer.git

RUN useradd -rm -d /home/ubuntu -s /bin/bash -g root -G sudo -u 1001 ubuntu
USER ubuntu
WORKDIR /repo
