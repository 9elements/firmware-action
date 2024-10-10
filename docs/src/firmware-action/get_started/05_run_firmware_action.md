# Run firmware-action locally

```bash
./firmware-action build --config=firmware-action.json --target=coreboot-example
```

`firmware-action` will firstly download `registry.dagger.io/engine` container needed for dagger and start it.

Then it will proceed to download `coreboot` container (specified by `sdk_url` in JSON config), copy into it specified files and then start compilation.

If compilation is successful, a new directory `output-coreboot/` will be created (as specified by `output_dir` in JSON config) which will contain files (specified by `container_output_files` in JSON config) and possibly also directories (specified by `container_output_dirs` in JSON config).

```admonish info
`container_output_dirs` and `container_output_files` are lists of directories and files to be extracted from the container once compilation finished successfully.

These are then placed into `output_dir`.
```

