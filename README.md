# Firmware-Action

[![Lint](https://github.com/9elements/firmware-action/actions/workflows/lint.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/lint.yml)
[![pytest](https://github.com/9elements/firmware-action/actions/workflows/pytest.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/pytest.yml)
[![dagger](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml)
[![go test](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml)

Firmware-Action is a tool to simplify building firmware. Think of it as `Makefile` or `Taskfile` but specifically for firmware. The tool it-self is written entirely in [Golang](https://go.dev/).

Motivation behind the creation is to unify building of firmware across development environments. The goal of firmware action is to run on your local machine but also in your CI/CD pipeline, with the same configuration producing the same output.

There is also an independent python tool to prepare Docker containers to be used with Firmware-Action. These are hosted on GitHub and are freely available (no need to build any Docker containers yourself).

There is also a GitHub action integration allowing you to use Firmware-Action in your GitHub CI/CD.

At the moment Firmware-Action supports:
- [coreboot](https://coreboot.org/)
- [linux](https://www.kernel.org/)
- [tianocore / edk2](https://www.tianocore.org/)
- firmware stitching (populating IFD regions with binaries)
- [u-root](https://github.com/u-root/u-root)

This list should expand in the future (see [issue 56](https://github.com/9elements/firmware-action/issues/56)).

Firmware-Action is using [dagger](https://docs.dagger.io/) under the hood, which makes it a rather versatile tool. When Firmware-Action is used, it automatically downloads user-specified Docker containers in which it will attempt to build the firmware.

If your firmware consists of multiple components, such as `coreboot` with `linux` as the payload, you can simply define each as a module and define dependency between them. Each module is built separately, but can use the output of another module as input. In the `coreboot` + `linux` example, you can call Firmware-Action to build `coreboot` recursively, which will also build `linux` due to the dependency definition. This way, you can build complex stacks of firmware in a single call.

[Documentation is hosted in pages](https://9elements.github.io/firmware-action/).

