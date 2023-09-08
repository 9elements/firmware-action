# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=buildpack-deps:focal
ARG INTERMEDIATE_IMAGE=base

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG EDK2_VERSION=edk2-stable202008
ARG EDK2_PLATFORM_COMMIT=7093026e924291b9e9f6f5211a052c6a761a4704
ARG EDK2_NON_OSI_COMMIT=4f88718028316aee31cb577f7127d5706255722d
ARG GCC_VERSION=9
ENV GCC_VERSION=$GCC_VERSION
ARG PYTHON_PACKAGES=python-is-python3 python2 python3
ARG GCC_CROSS_COMPILER_PACKAGES=gcc-${GCC_VERSION}-aarch64-linux-gnu gcc-${GCC_VERSION}-arm-linux-gnueabi gcc-${GCC_VERSION}-i686-linux-gnu

ARG CONTEXT=edk2

ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_EDK2_VERSION=$EDK2_VERSION

ENV DEBIAN_FRONTEND=noninteractive
# Directory for coreboot toolchain, MEAnalyser etc.
ENV TOOLSDIR=/tools
# EDK2 required env vars
ENV WORKSPACE_CORE=$TOOLSDIR/Edk2
ENV WORKSPACE=$TOOLSDIR/Edk2

RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        ${PYTHON_PACKAGES} \
        acpica-tools \
        bsdmainutils \
        g++-${GCC_VERSION} \
        gcc-${GCC_VERSION} \
        ${GCC_CROSS_COMPILER_PACKAGES} \
        git \
        iucode-tool \
        nasm \
        nodejs \
        openssh-client \
        qemu-system-x86 \
        uuid-dev \
        wine-stable \
        wine32 \
        wine64 \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}" && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCC_VERSION} 100


#=============
# Dependencies needed for universal payload
FROM base AS universalpayload

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        clang \
        llvm \
        lld \
        gcc-multilib \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*


#=============
# "toolchain" stage to build stuff
#   should either build off "base" or "universalpayload"
ARG INTERMEDIATE_IMAGE
FROM ${INTERMEDIATE_IMAGE} AS toolchain

WORKDIR $TOOLSDIR
RUN git clone --recurse-submodules https://github.com/tianocore/edk2-platforms.git Edk2Platforms && \
    cd Edk2Platforms && \
    git checkout "${EDK2_PLATFORM_COMMIT}"
RUN git clone --recurse-submodules https://github.com/tianocore/edk2-non-osi.git Edk2NonOsi && \
    cd Edk2NonOsi && \
    git checkout "${EDK2_NON_OSI_COMMIT}"
RUN git clone --recurse-submodules https://github.com/tianocore/edk2.git --depth 2 --branch "${EDK2_VERSION}" Edk2 && \
    cd Edk2 && \
    make -C BaseTools/ -j "$(nproc)"


#=============
# "final" stage is the actual product with everything included
#   should either build off "base" or "universalpayload"
ARG INTERMEDIATE_IMAGE
FROM ${INTERMEDIATE_IMAGE} AS final

COPY --from=toolchain $TOOLSDIR $TOOLSDIR
