# Get firmware-action

Firstly, you will need to [install and setup Docker](https://docs.docker.com/engine/install/).

Then you can get firmware-action multiple ways:

## Build from source
Git clone and build, we use [Taskfile](https://taskfile.dev/) as build system, but you can go with just `go build`.
```
git clone https://github.com/9elements/firmware-action.git
cd firmware-action
task build-go-binary
```

## Download executable
Download pre-compiled executable from [releases](https://github.com/9elements/firmware-action/releases).

## Arch Linux
There is [AUR package](https://aur.archlinux.org/packages/firmware-action) available.

## go install
```
go install -v github.com/9elements/firmware-action/cmd/firmware-action@latest
```
