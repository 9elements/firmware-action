# Run firmware-action locally

```bash
./firmware-action build --config=firmware-action.json --target=coreboot-example
```

`firmware-action` will firstly download `registry.dagger.io/engine` container needed for dagger and start it.

It will then proceed to download a `coreboot` container {{footnote: The used container is specified by `sdk_url` in the firmware-action configuration file.}}, copy into it the specified files and then start compilation.

If compilation is successful, a new directory `output-coreboot/` will be created {{footnote: Output directory is specified by `output_dir` in firmware-action configuration file.}} which will contain files {{footnote: Output files are specified by `container_output_files` in firmware-action configuration file.}} and possibly also directories {{footnote: Directories to output are specified by `container_output_dirs` in firmware-action configuration file.}}.

Your working directory should look something like this:
```
.
|-- coreboot/
|   `-- ...
|-- firmware-action.json
|-- output-coreboot/
|   |-- coreboot.rom
|   `-- defconfig
`-- seabios_defconfig
```


```admonish info
`container_output_dirs` and `container_output_files` are lists of directories and files to be extracted from the container once compilation finished successfully.

These are then placed into `output_dir`.
```

