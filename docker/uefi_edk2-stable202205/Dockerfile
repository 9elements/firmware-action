# syntax=docker/dockerfile:experimental
FROM buildpack-deps:jammy

ARG VERSION=edk2-stable202205
ARG GCCVERSION=9
ARG PLATFORMVERSION=f427247a8d415c2d514ee49c4d0cde94a9b8ea89
ARG NONOSIVERION=6996a45

ENV DEBIAN_FRONTEND=noninteractive
ENV WORKSPACE_CORE=/repo/Edk2

# install tools
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y openssh-client git gcc-${GCCVERSION} g++-${GCCVERSION} \
        python-is-python3 python2 python3 \
        uuid-dev acpica-tools nano vim nasm wine-stable wine32 wine64 bsdmainutils less \
        iucode-tool && \
    mkdir -p /repo/ && \
    cd /repo/; \
    git clone --recurse-submodules https://github.com/tianocore/edk2.git --depth 2 --branch ${VERSION} Edk2; \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCCVERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCCVERSION} 100 && \
    git clone --recurse-submodules https://github.com/tianocore/edk2-platforms.git Edk2Platforms && \
    git clone --recurse-submodules https://github.com/tianocore/edk2-non-osi.git Edk2NonOsi && \
    cd /repo/Edk2Platforms; git checkout ${PLATFORMVERSION} && \
    cd /repo/Edk2NonOsi; git checkout ${NONOSIVERSION} && \
    cd /repo/Edk2; make -C BaseTools/ -j $(nproc) && \
    winecfg

WORKDIR /repo
