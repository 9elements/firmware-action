# Firmware-Action

[![Lint](https://github.com/9elements/firmware-action/actions/workflows/lint.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/lint.yml)
[![building containers](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml)
[![go test](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/9elements/firmware-action/cmd/firmware-action)](https://goreportcard.com/report/github.com/9elements/firmware-action/cmd/firmware-action)
[![codecov](https://codecov.io/gh/9elements/firmware-action/graph/badge.svg?token=FDWKRAJIJ0)](https://codecov.io/gh/9elements/firmware-action)


## Description

`Firmware-action` is a tool to simplify building firmware. Think of it as `Makefile` or `Taskfile` but specifically for firmware. The tool it-self is written entirely in [Golang](https://go.dev/).

Motivation behind the creation is to unify building of firmware across development environments. The goal of `firmware-action` is to run on your local machine but also in your CI/CD pipeline, with the same configuration producing the same output.

There is also an independent python tool to prepare Docker containers to be used with `firmware-action`. These are hosted on GitHub and are freely available (no need to build any Docker containers yourself).

There is also a GitHub action integration allowing you to use `firmware-action` in your GitHub CI/CD.

At the moment `firmware-action` has modules to build:
- [coreboot](https://coreboot.org/)
- [linux](https://www.kernel.org/)
- [tianocore / edk2](https://www.tianocore.org/)
- firmware stitching (populating IFD regions with binaries)
- [u-root](https://github.com/u-root/u-root)
- [u-boot](https://docs.u-boot.org/en/latest/index.html)
- universal module (run arbitrary commands in arbitrary Docker container)

This list should expand in the future (see [issue 56](https://github.com/9elements/firmware-action/issues/56)).

`Firmware-action` is using [dagger](https://docs.dagger.io/) under the hood, which makes it a rather versatile tool. When `firmware-action` is used, it automatically downloads user-specified Docker containers in which it will attempt to build the firmware.

If your firmware consists of multiple components, such as `coreboot` with `linux` as the payload, you can simply define each as a module and define dependency between them. Each module is built separately, but can use the output of another module as input. In the `coreboot` + `linux` example, you can call `firmware-action` to build `coreboot` recursively, which will also build `linux` due to the dependency definition. This way, you can build complex stacks of firmware in a single call.

[Documentation is hosted in pages](https://9elements.github.io/firmware-action/).

There is a standalone repository with usage examples at [firmware-action-example](https://github.com/9elements/firmware-action-example).

Pre-compiled coreboot toolchains for docker containers are stored in separate repository [firmware-action-toolchains](https://github.com/9elements/firmware-action-toolchains).


## Containers

We maintain multiple containers that can be freely used with `firmware-action`. They are hosted both in [GitHub registry](https://github.com/orgs/9elements/packages?repo_name=firmware-action) as well as in [DockerHub](https://hub.docker.com/u/9elementscyberops).

DockerHub registry contains only releases, while GitHub registry contains also containers build on `main` branch.

Here is a list of all containers:

| Container            | Maintained | Note                                       |
| -------------------- | ---------- | ------------------------------------------ |
| coreboot_4.19        | [x]        |                                            |
| coreboot_4.20        | [ ]        | discontinued in favor of `4.20.1`          |
| coreboot_4.20.1      | [x]        |                                            |
| coreboot_4.21        | [x]        |                                            |
| coreboot_4.22        | [ ]        | discontinued in favor of `4.22.1`          |
| coreboot_4.22.01     | [x]        |                                            |
| coreboot_24.02       | [ ]        | discontinued in favor of `24.02.01`        |
| coreboot_24.02.01    | [x]        |                                            |
| coreboot_24.05       | [x]        |                                            |
| coreboot_25.03       | [x]        |                                            |
| coreboot_25.06       | [x]        |                                            |
| udk2017              | [x]        |                                            |
| edk2-stable202008    | [x]        |                                            |
| edk2-stable202105    | [x]        |                                            |
| edk2-stable202111    | [x]        |                                            |
| edk2-stable202205    | [x]        |                                            |
| edk2-stable202208    | [x]        |                                            |
| edk2-stable202211    | [x]        |                                            |
| edk2-stable202302    | [x]        |                                            |
| edk2-stable202305    | [x]        |                                            |
| edk2-stable202308    | [x]        |                                            |
| edk2-stable202311    | [x]        |                                            |
| edk2-stable202402    | [x]        |                                            |
| edk2-stable202405    | [x]        |                                            |
| edk2-stable202408    | [ ]        | discontinued in favor of `stable202408.01` |
| edk2-stable202408.01 | [x]        |                                            |
| edk2-stable202411    | [x]        |                                            |
| edk2-stable202502    | [x]        |                                            |
| edk2-stable202505    | [x]        |                                            |
| edk2-stable202508    | [x]        |                                            |
| linux_6.1.111        | [ ]        | discontinued in favor of `linux_6.1`       |
| linux_6.1.45         | [ ]        | discontinued in favor of `linux_6.1`       |
| linux_6.1            | [x]        |                                            |
| linux_6.6.52         | [ ]        | discontinued in favor of `linux_6.6`       |
| linux_6.6            | [x]        |                                            |
| linux_6.9.9          | [ ]        | discontinued because not LTS               |
| linux_6.11           | [ ]        | discontinued because not LTS               |
| linux_6.12           | [x]        |                                            |
| uroot_0.14.0         | [x]        |                                            |
| uroot_0.15.0         | [x]        |                                            |
| uboot_2025.01        | [x]        |                                            |
| uboot_2025.04        | [x]        |                                            |
| uboot_2025.07        | [x]        |                                            |


## Legacy containers

These were created by hand long time ago and since then have been replaced.

| Container                                                                         | Maintained | Note         |
| --------------------------------------------------------------------------------- | ---------- | ------------ |
| [coreboot](https://github.com/orgs/9elements/packages/container/package/coreboot) | [ ]        | discontinued |
| [uefi](https://github.com/orgs/9elements/packages/container/package/uefi)         | [ ]        | discontinued |
