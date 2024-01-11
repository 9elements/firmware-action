# Compile Firmware

[![Lint](https://github.com/9elements/firmware-action/actions/workflows/lint.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/lint.yml)
[![pytest](https://github.com/9elements/firmware-action/actions/workflows/pytest.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/pytest.yml)
[![dagger](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/docker-build-and-test.yml)
[![go test](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml/badge.svg)](https://github.com/9elements/firmware-action/actions/workflows/go-test.yml)

This repository contains tools to simplify builds of firmware stacks.

At the moment it supports:
- [coreboot](https://coreboot.org/)
- [linux](https://www.kernel.org/)
- [tianocore / edk2](https://www.tianocore.org/)

This list should expand in the future (see [issue 56](https://github.com/9elements/firmware-action/issues/56)).

Motivation behind the creation is to unify building of firmware across development. The same code and configuration should run in CI/CD pipeline just as well as on your local machine.

Initially it was meant only as GitHub-specific action, but it should be universal thanks to [dagger](https://docs.dagger.io/).

