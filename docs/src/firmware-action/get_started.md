# Get started

This guide will provide instructions step by step on how to get started with firmware-action, and it will demonstrate the use on coreboot example.

In this guide we will:
- start a new repository
    - in this guide it will be hosted in [GitHub](https://github.com/)
- we will build a simple coreboot for [QEMU](https://wiki.archlinux.org/title/QEMU)
- we will be able to build coreboot in GitHub action and locally

The code from this example is available in [firmware-action-example](https://github.com/9elements/firmware-action-example).


## Prerequisites

- installed [Docker](https://wiki.archlinux.org/title/Docker)
- installed git


## Start a new git repository

Start a new repository in GitHub and the clone it.


## Add coreboot as submodule

~~~admonish tip
Add git submodule with:
```
git submodule add <repo> <path>
```
~~~

Add [coreboot repository](https://review.coreboot.org/admin/repos/coreboot,general) as a submodule:
```bash
git submodule add --depth=1 "https://review.coreboot.org/coreboot" coreboot
```

```bash
git submodule update --init
```

Checkout a release tag, for example `4.19` (it is a bit older release, but should suffice for demonstration)
```bash
( cd coreboot; git fetch origin tag "4.19"; git checkout "4.19" )
```


~~~admonish warning
Recursively initializing all submodules in coreboot will take a moment.
~~~

```bash
git submodule update --init --recursive
```


## Create a coreboot configuration file

Either follow a [coreboot guide](https://doc.coreboot.org/tutorial/part1.html#step-5-configure-the-build) to get a base-bones-basic configuration, or just copy-paste this text into `seabios_defconfig` file.

~~~admonish example title="seabios_defconfig"
```properties
{{#include ../firmware-action-example/seabios_defconfig}}
```
~~~


## Create a JSON configuration file

This configuration file is for firmware-action, so that it knows what to do and where to find things. Let's call it `firmware-action.json`.

~~~admonish example title="firmware-action.json"
```json
{{#include ../firmware-action-example/firmware-action.json}}
```
~~~

~~~admonish info
Field `repo_path` is pointing to the location of out coreboot submodule.
~~~

~~~admonish info
Field `defconfig_path` is pointing to the location of coreboot's configuration file.
~~~


## Get firmware-action

Either [clone the repository and build the executable yourself](./usage_local.md), or just download pre-compiled executable from [releases](https://github.com/9elements/firmware-action/releases).


## Run firmware-action locally

```
./firmware-action build --config=firmware-action.json --target=coreboot-example
```

`firmware-action` will firstly download `registry.dagger.io/engine` container needed for dagger and start it.

Then it will proceed to download `coreboot` container (specified by `sdk_url` in JSON config), copy into it specified files and then start compilation.

If compilation is successful, a new directory `output-coreboot/` will be created (as specified by `output_dir` in JSON config) which will contain files (specified by `container_output_files` in JSON config) and possibly also directories (specified by `container_output_dirs` in JSON config).

~~~admonish info
`container_output_dirs` and `container_output_files` are lists of directories and files to be extracted from the container once compilation finished successfully.

These are then placed into `output_dir`.
~~~


## Run firmware-action in GitHub CI

Now that we have `firmware-action` working on local system. Let's set up CI.

~~~admonish example title=".github/workflows/example.yml"
```yaml
{{#include ../firmware-action-example/.github/workflows/example.yml}}
```
~~~

And that is it.

