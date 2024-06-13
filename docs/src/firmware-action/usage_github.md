# Github CI

You can use `firmware-action` as any other action.

```admonish example
~~~yaml
name: Firmware example build
jobs:
  firmware_build:
    runs-on: ubuntu-latest
    steps:
      - name: Build coreboot with firmware-action
        uses: 9elements/firmware-action@main
        with:
          config: '<path to firmware-action JSON config>'
          target: '<name of the target from JSON config>'
          recursive: 'false'
~~~
```


## Parametric builds with environment variables

To take advantage of matrix builds in GitHub, it is possible to use environment variables inside the JSON configuration file.

```admonish example
For example let's make `COREBOOT_VERSION` environment variable which will hold version of coreboot.

JSON would look like this:
~~~json
...
"sdk_url": "ghcr.io/9elements/firmware-action/coreboot_${COREBOOT_VERSION}:main",
...
"defconfig_path": "tests/coreboot_${COREBOOT_VERSION}/seabios.defconfig",
...
~~~

YAML would look like this:
~~~yaml
name: Firmware example build
jobs:
  firmware_build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        coreboot_version: ["4.19", "4.20"]
    steps:
      - name: Build coreboot with firmware-action
        uses: 9elements/firmware-action@main
        with:
          config: '<path to firmware-action JSON config>'
          target: '<name of the target from JSON config>'
          recursive: 'false'
        env:
          COREBOOT_VERSION: ${{ matrix.coreboot_version }}
~~~
```


## Examples

In our repository we have multiple examples (even though rather simple ones) defined in [.github/workflows/example.yml](https://github.com/9elements/firmware-action/blob/main/.github/workflows/example.yml).

```admonish example collapsible=true title="Coreboot"
`.github/workflows/example.yml`:
~~~yaml
{{#include ../../../.github/workflows/example.yml:example_build_coreboot}}
~~~

`tests/example_config__coreboot.json`:
~~~json
{{#include ../../../tests/example_config__coreboot.json}}
~~~
```

```admonish example collapsible=true title="Linux Kernel"
`.github/workflows/example.yml`:
~~~yaml
{{#include ../../../.github/workflows/example.yml:example_build_linux_kernel}}
~~~

`tests/example_config__linux.json`:
~~~json
{{#include ../../../tests/example_config__linux.json}}
~~~
```

```admonish example collapsible=true title="Edk2"
`.github/workflows/example.yml`:
~~~yaml
{{#include ../../../.github/workflows/example.yml:example_build_edk2}}
~~~

`tests/example_config__edk2.json`:
~~~json
{{#include ../../../tests/example_config__edk2.json}}
~~~
```

```admonish example collapsible=true title="Firmware Stitching"
`.github/workflows/example.yml`:
~~~yaml
{{#include ../../../.github/workflows/example.yml:example_build_stitch}}
~~~

`tests/example_config__firmware_stitching.json`:
~~~json
{{#include ../../../tests/example_config__firmware_stitching.json}}
~~~
```

