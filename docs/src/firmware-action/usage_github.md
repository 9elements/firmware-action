# Github CI

You can use `firmware-action` as any other action. The action automatically handles artifact uploading and caching to improve build times.


## Artifacts and Outputs

The action automatically uploads build artifacts and provides the following outputs:
- `artifact-name`: Name of the uploaded artifact

Internally, `firmware-action` is using [upload-artifact](https://github.com/actions/upload-artifact) action and passes on it's outputs too:
- `artifact-id`: GitHub ID of the artifact (useful for REST API)
- `artifact-url`: Direct download URL for the artifact (requires GitHub login)
- `artifact-digest`: SHA-256 digest of the artifact

```admonish example
You can use these outputs in subsequent steps:
~~~yaml
steps:
  - name: Build firmware
    id: firmware
    uses: 9elements/firmware-action@main
    with:
      config: 'config.json'
      target: 'my-target'

  - name: Use artifact info
    run: |
      echo "Artifact name: ${{ steps.firmware.outputs.artifact-name }}"
      echo "Download from: ${{ steps.firmware.outputs.artifact-url }}"
~~~
```


## Build Caching

The action automatically caches build artifacts between runs to speed up builds. The cache is:
- Keyed by the config file contents and commit SHA
- Restored at the start of each run
- Saved after each run, even if the build fails

No configuration is needed - caching works out of the box.


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
