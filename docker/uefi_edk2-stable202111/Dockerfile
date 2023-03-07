# syntax=docker/dockerfile:experimental
FROM buildpack-deps:focal

ARG VERSION=edk2-stable202111
ARG GCCVERSION=9
ARG PLATFORMVERSION=41dacdf4fb36b350bd199adeb9036b7d44b3d243
ARG NONOSIVERSION=eef5e03

ENV DEBIAN_FRONTEND=noninteractive
ENV WORKSPACE_CORE=/repo/Edk2

# install tools
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y openssh-client git gcc-${GCCVERSION} g++-${GCCVERSION} \
      python-is-python3 python2 python3 \
      uuid-dev acpica-tools nano vim nasm wine-stable wine32 wine64 bsdmainutils less \
      clang llvm lld gcc-multilib \
      nodejs \
      iucode-tool && apt-get clean && \
    update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-${GCCVERSION} 100 && \
    update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-${GCCVERSION} 100 && \
    mkdir -p /repo

WORKDIR /repo
RUN git clone --recurse-submodules https://github.com/tianocore/edk2.git --depth 2 --branch ${VERSION} Edk2; \
    git clone --recurse-submodules https://github.com/tianocore/edk2-platforms.git Edk2Platforms && \
    git clone --recurse-submodules https://github.com/tianocore/edk2-non-osi.git Edk2NonOsi

WORKDIR /repo/Edk2Platforms
RUN git checkout ${PLATFORMVERSION}

WORKDIR /repo/Edk2NonOsi
RUN git checkout ${NONOSIVERSION}

WORKDIR /repo/Edk2
RUN make -C BaseTools/ -j $(nproc) && winecfg

WORKDIR /repo
