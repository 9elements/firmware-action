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

To take advantage of matrix builds, it is possible to use environment variables inside the JSON configuration file.

```admonish example
For example let's make `COREBOOT_VERSION` environment variable which will hold verion of coreboot.

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

