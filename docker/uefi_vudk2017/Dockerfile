# syntax=docker/dockerfile:experimental
FROM buildpack-deps:bionic

ARG VERSION=vUDK2017
ARG GCCVERSION=4.8
# NOTE: Platforms wasn't used in 2017
ARG PLATFORMVERSION=5b849a6397
# NOTE: NonOSI wasn't used in 2017
ARG NONOSIVERSION=6633222

ENV DEBIAN_FRONTEND=noninteractive
ENV WORKSPACE_CORE=/repo/Edk2

# install tools
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y openssh-client git gcc-${GCCVERSION} g++-${GCCVERSION} \
        python \
        uuid-dev acpica-tools nano vim nasm wine-stable wine32 wine64 bsdmainutils less \
        iucode-tool && \
    mkdir -p /repo/ && \
    cd /repo/; \
    git clone --recurse-submodules https://github.com/tianocore/edk2.git --depth 2 --branch ${VERSION} Edk2; \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCCVERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCCVERSION} 100  && \
    git clone --recurse-submodules https://github.com/tianocore/edk2-platforms.git Edk2Platforms && \
    git clone --recurse-submodules https://github.com/tianocore/edk2-non-osi.git Edk2NonOsi && \
    cd /repo/Edk2Platforms; git checkout ${PLATFORMVERSION} && \
    cd /repo/Edk2NonOsi; git checkout ${NONOSIVERSION} && \
    cd /repo/Edk2; make -C BaseTools/ && \
    winecfg

WORKDIR /repo
