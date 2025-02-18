# Configuration

```admonish example collapsible=true title="Example of JSON configuration file"
~~~json
{{#include ../../../tests/example_config.json}}
~~~
```

```admonish tip
Multiple configuration files can be supplied to firmware-action. Dependencies also work across files.

~~~
firmware-action build --config=config-01.json --config=config-02.json ...
~~~

Beware that modules with identical names are permitted, as long as they are not in the same configuration file.

`firmware-action` processes the files in order in which they were supplied and in case of name-collision, the configuration in last file takes precedence.
```

The configuration is split by type (`coreboot`, `linux`, `edk2`, ...).

In each type can be any number of modules.

Each module has a name, which can be anything as long as it is unique (unique string across all modules of all types). In the example above there are 3 modules (`coreboot-example`, `linux-example`, `edk2-example`).

The configuration above can be simplified to this:
```
/
├── coreboot/
│   └── coreboot-example
├── edk2/
│   └── edk2-example
├── firmware_stitching/
│   └── stitching-example
└── linux/
    └── linux-example
```

Not all types must be present or defined. If you are building coreboot and coreboot only, you can have only coreboot present.
```
/
└── coreboot/
    └── coreboot_example
```

You can have multiple modules of each type, as long as their names are unique.
```
/
├── coreboot/
│   ├── coreboot_example
│   ├── coreboot_A
│   └── my_little_firmware
├── linux/
│   ├── linux_example
│   ├── linux_B
│   ├── asdf
│   └── asdf2
└── edk2/
    ├── edk2_example
    └── edk2_C
```


## Modules

Each module has sections:
- `depends`
- `common`
- `specific`

```go
{{#include ../../../cmd/firmware-action/recipes/coreboot.go:CorebootOpts}}
```

`common` & `specific` are identical in function. There is no real difference between these two. They are split to simplify the code. They define things like path to source code, version and source of SDK to use, and so on.

`depends` on the other hand allows you to specify dependency (or relation) between modules. For example your `coreboot` uses `edk2` as payload. So you can specify this dependency by listing name of the `edk2` module in `depends` of your `coreboot` module.

```json
{
  "coreboot": {
    "coreboot-example": {
      "depends": ["edk2-example"],
      ...
    }
  },
  "edk2": {
    "edk2-example": {
      "depends": null,
      ...
    }
  }
}
```

With such configuration, you can then run `firmware-action` recursively, and it will build all of the modules in proper order.
```
./firmware-action build --config=./my-config.json --target=coreboot-example --recursive
```
In this case `firmware-action` would build `edk2-example` first and then `coreboot-example`.

```admonish tip
By changing inputs and outputs, you can then feed output of one module into input of another module.

This way you can build the entire firmware stack in single step.
```


## Common and Specific

To explain each and every entry in the configuration, here are snippets of the source code with comments.

```admonish info
In the code below, the tag `json` (for example `json:"sdk_url"`) specifies what the field is called in JSON file.

Tag `validate:"required"`, it means that the field is required and must not be empty. Empty required field will fail validation and terminate program with error.

Tag `validate:"dirpath"` means that field must contain a valid path to a directory. It is not necessary for the path or directory to exists, but must be a valid path. Be warned - that means that the string must end with `/`. For example `output-coreboot/`.

Tag `validate:"filepath"` means that the field must contain a valid path to a file. It is not necessary for the file to exist.

For more tails see [go-playground/validator](https://github.com/go-playground/validator) package.
```

### Common
```go
{{#include ../../../cmd/firmware-action/recipes/config.go:CommonOpts}}
```

### Specific / coreboot
```go
{{#include ../../../cmd/firmware-action/recipes/coreboot.go:CorebootOpts}}
{{#include ../../../cmd/firmware-action/recipes/coreboot.go:CorebootBlobs}}
```

### Specific / Linux
```go
{{#include ../../../cmd/firmware-action/recipes/linux.go:LinuxOpts}}
{{#include ../../../cmd/firmware-action/recipes/linux.go:LinuxSpecific}}
```

### Specific / Edk2
```go
{{#include ../../../cmd/firmware-action/recipes/edk2.go:Edk2Opts}}
{{#include ../../../cmd/firmware-action/recipes/edk2.go:Edk2Specific}}
```

### Specific / Firmware stitching
```go
{{#include ../../../cmd/firmware-action/recipes/stitching.go:FirmwareStitchingOpts}}
{{#include ../../../cmd/firmware-action/recipes/stitching.go:IfdtoolEntry}}
```

### Specific / u-root
```go
{{#include ../../../cmd/firmware-action/recipes/uroot.go:URootOpts}}
{{#include ../../../cmd/firmware-action/recipes/uroot.go:URootSpecific}}
```

### Specific / Universal module
```go
{{#include ../../../cmd/firmware-action/recipes/universal.go:UniversalOpts}}
{{#include ../../../cmd/firmware-action/recipes/universal.go:UniversalSpecific}}
```

### Specific / u-boot module
```go
{{#include ../../../cmd/firmware-action/recipes/uboot.go:UBootOpts}}
```
