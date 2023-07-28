# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

#=============
# "base" stage with all needed build dependencies
FROM buildpack-deps:focal AS base

ARG VERSION=edk2-stable202008
ARG GCCVERSION=9
ARG PLATFORMVERSION=7093026e924291b9e9f6f5211a052c6a761a4704
ARG NONOSIVERSION=4f88718

ENV DEBIAN_FRONTEND=noninteractive
# Directory for coreboot toolchain, MEAnalyser etc.
ENV TOOLSDIR=/tools
# EDK2 required env vars
ENV WORKSPACE_CORE=$TOOLSDIR/Edk2
ENV WORKSPACE=$TOOLSDIR/Edk2

RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        acpica-tools \
        bsdmainutils \
        g++-${GCCVERSION} \
        gcc-${GCCVERSION} \
        git \
        iucode-tool \
        less \
        nano \
        nasm \
        openssh-client \
        python-is-python3 \
        python2 \
        python3 \
        qemu-system-x86 \
        uuid-dev \
        vim \
        wine-stable \
        wine32 \
        wine64 \
    && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}"


#=============
# "toolchain" stage to build stuff
FROM base AS toolchain

WORKDIR $TOOLSDIR
RUN update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCCVERSION} ${GCCVERSION} && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCCVERSION} ${GCCVERSION}
RUN git clone --recurse-submodules https://github.com/tianocore/edk2-platforms.git Edk2Platforms && \
    cd Edk2Platforms && \
    git checkout ${PLATFORMVERSION}
RUN git clone --recurse-submodules https://github.com/tianocore/edk2-non-osi.git Edk2NonOsi && \
    cd Edk2NonOsi && \
    git checkout ${NONOSIVERSION}
RUN git clone --recurse-submodules https://github.com/tianocore/edk2.git --depth 2 --branch ${VERSION} Edk2 && \
    cd Edk2 && \
    make -C BaseTools/ -j $(nproc)
RUN winecfg


#=============
# "final" stage is the actual product with everything included
FROM base AS final

COPY --from=toolchain $TOOLSDIR $TOOLSDIR
WORKDIR /repo
