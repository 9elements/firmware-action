# This is multi-stage Dockerfile
#   Docs: https://docs.docker.com/build/building/multi-stage/

ARG SOURCE_IMAGE=buildpack-deps:focal
ARG INTERMEDIATE_IMAGE=base
ARG TARGETARCH=amd64

#=============
# "base" stage with all needed build dependencies
FROM ${SOURCE_IMAGE} AS base

ARG TARGETARCH

ARG EDK2_VERSION=edk2-stable202008
ARG EDK2_VERSION_COMMIT=06dc822d045c2bb42e497487935485302486e151
ARG EDK2_PLATFORM_COMMIT=7093026e924291b9e9f6f5211a052c6a761a4704
ARG EDK2_NON_OSI_COMMIT=4f88718028316aee31cb577f7127d5706255722d

ARG GCC_VERSION=9
ENV GCC_VERSION=$GCC_VERSION
ARG USE_GCC_VERSION=5
ENV USE_GCC_VERSION=$USE_GCC_VERSION

ARG PYTHON_PACKAGES=python-is-python3 python2 python3
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
        git \
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
        less \
        nano \
        openssh-server \
        screen \
        tmux \
        tree \
        vim \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p "${TOOLSDIR}" && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCC_VERSION} 100 && \
    update-alternatives --install /usr/local/bin/python python /usr/bin/python3 100 && \
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
    git checkout "${EDK2_VERSION_COMMIT}" && \
    git submodule update --init --recursive && \
    if [ "${TARGETARCH}" = 'amd64' ]; then \
        make -C BaseTools/ -j "$(nproc)"; \
    fi;


#=============
# "final" stage is the actual product with everything included
#   should either build off "base" or "universalpayload"
ARG INTERMEDIATE_IMAGE
FROM ${INTERMEDIATE_IMAGE} AS final

COPY --from=toolchain $TOOLSDIR $TOOLSDIR

# Prepare SSH for interactive debugging
RUN mkdir -p /run/sshd && \
    sed -E "s/#PermitRootLogin\b.*/PermitRootLogin yes/g" -i /etc/ssh/sshd_config && \
    sed -E "s/#PasswordAuthentication\b.*/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config

