# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=buildpack-deps:focal
ARG INTERMEDIATE_IMAGE=base
ARG TARGETARCH=amd64
ARG NODESOURCE=https://deb.nodesource.com/setup_24.x

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH
ARG NODESOURCE

ARG EDK2_VERSION=edk2-stable202008
ENV EDK2_VERSION=$EDK2_VERSION
ARG EDK2_VERSION_COMMIT=06dc822d045c2bb42e497487935485302486e151
ENV EDK2_VERSION_COMMIT=$EDK2_VERSION_COMMIT
ARG EDK2_PLATFORM_COMMIT=7093026e924291b9e9f6f5211a052c6a761a4704
ENV EDK2_PLATFORM_COMMIT=$EDK2_PLATFORM_COMMIT
ARG EDK2_NON_OSI_COMMIT=4f88718028316aee31cb577f7127d5706255722d
ENV EDK2_NON_OSI_COMMIT=$EDK2_NON_OSI_COMMIT

ARG GCC_VERSION=9
ENV GCC_VERSION=$GCC_VERSION
ARG USE_GCC_VERSION=5
ENV USE_GCC_VERSION=$USE_GCC_VERSION

ARG PYTHON_PACKAGES=python-is-python3 python2 python3
ARG PYTHON_VERSION=python3
ARG GCC_CROSS_COMPILER_PACKAGES_ARM=gcc-${GCC_VERSION}-aarch64-linux-gnu gcc-${GCC_VERSION}-arm-linux-gnueabi
ARG GCC_CROSS_COMPILER_PACKAGES_X86=gcc-${GCC_VERSION}-i686-linux-gnu

ARG CONTEXT=edk2

ENV VERIFICATION_TEST=./tests/test_${CONTEXT}.sh
ENV VERIFICATION_TEST_EDK2_VERSION=$EDK2_VERSION

ENV DEBIAN_FRONTEND=noninteractive
# Directory for coreboot toolchain, MEAnalyser etc.
ENV TOOLSDIR=/tools
# EDK2 required env vars
ENV WORKSPACE_CORE=$TOOLSDIR/Edk2
ENV WORKSPACE=$TOOLSDIR/Edk2

RUN if [ "${TARGETARCH}" = 'amd64' ]; then \
        dpkg --add-architecture i386; \
    fi; \
    wget --quiet -O nodesource_setup.sh "${NODESOURCE}" && \
        chmod +x nodesource_setup.sh && \
        ./nodesource_setup.sh && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        ${PYTHON_PACKAGES} \
        acpica-tools \
        bison \
        bsdmainutils \
        build-essential \
        flex \
        g++-${GCC_VERSION} \
        gcc-${GCC_VERSION} \
        imagemagick \
        nasm \
        nodejs \
        openssh-client \
        qemu-system-x86 \
        uuid-dev \
        wine-stable \
        wine64 \
    && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        apt-get install -y --no-install-recommends \
            ${GCC_CROSS_COMPILER_PACKAGES_ARM} \
            iucode-tool \
            wine32 \
        ; \
    fi; \
    if [ "${TARGETARCH}" = 'arm64' ]; then \
        apt-get install -y --no-install-recommends \
            ${GCC_CROSS_COMPILER_PACKAGES_X86} \
        ; \
    fi; \
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
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}" && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/local/bin/python python /usr/bin/${PYTHON_VERSION} 100 && \
    wget --quiet -O rustup.sh https://sh.rustup.rs && \
        chmod +x ./rustup.sh && \
        ./rustup.sh -y --profile minimal


#=============
# Dependencies needed for universal payload
FROM base AS universalpayload
ARG TARGETARCH

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        clang \
        llvm \
        lld \
    && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        apt-get install -y --no-install-recommends \
            gcc-multilib \
        ; \
    fi; \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*


#=============
# "toolchain" stage to build stuff
#   should either build off "base" or "universalpayload"
ARG INTERMEDIATE_IMAGE
FROM ${INTERMEDIATE_IMAGE} AS toolchain

ARG TARGETARCH

# workaround to get shallow clone of specific commit
#   https://stackoverflow.com/questions/31278902/how-to-shallow-clone-a-specific-commit-with-depth-1/43136160#43136160
WORKDIR $TOOLSDIR
RUN mkdir Edk2Platforms && \
    cd Edk2Platforms && \
    git init && \
    git remote add origin https://github.com/tianocore/edk2-platforms.git && \
    git fetch --depth 1 origin "${EDK2_PLATFORM_COMMIT}" && \
    git checkout "${EDK2_PLATFORM_COMMIT}" && \
    git submodule update --init --recursive
RUN mkdir Edk2NonOsi && \
    cd Edk2NonOsi && \
    git init && \
    git remote add origin https://github.com/tianocore/edk2-non-osi.git && \
    git fetch --depth 1 origin "${EDK2_NON_OSI_COMMIT}" && \
    git checkout "${EDK2_NON_OSI_COMMIT}" && \
    git submodule update --init --recursive
RUN mkdir Edk2 && \
    cd Edk2 && \
    git init && \
    git remote add origin https://github.com/tianocore/edk2.git && \
    git fetch --depth 1 origin "${EDK2_VERSION_COMMIT}" && \
    git checkout "${EDK2_VERSION_COMMIT}"

# Back-porting patch 4dfdca63a9 which fixed submodules
COPY patches $TOOLSDIR/Edk2/.firmware-action-patches
ENV PATCH_FILE=.firmware-action-patches/${EDK2_VERSION}-zeex-subhook.patch
RUN cd Edk2 && \
    if [ -f "${PATCH_FILE}" ]; then \
        git apply "${PATCH_FILE}"; \
    fi && \
    if [ "${EDK2_VERSION}" = 'edk2-stable202008' ]; then \
        git submodule update --init --recursive; \
    else \
        git submodule update --init --recursive --depth 1; \
    fi && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        make -C BaseTools/ -j "$(nproc)"; \
    fi;

# For some reason edk2-stable202008 fails to clone all of the git submodules when shallow recursive init is used
# So to fix this we need to add IF ELSE statement to not use shallow clone on git submodules for edk2-stable202008
# Shame on you ekd2, shame on you


#=============
# "final" stage is the actual product with everything included
#   should either build off "base" or "universalpayload"
ARG INTERMEDIATE_IMAGE
FROM ${INTERMEDIATE_IMAGE} AS final

COPY --chown=root:root --chmod=755 scripts/* /bin/
COPY --from=toolchain $TOOLSDIR $TOOLSDIR

