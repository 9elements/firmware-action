# Configuration

```admonish example
Example of JSON configuration file:
~~~json
{{#include ../../../tests/example_config.json}}
~~~
```

The config are split by type (coreboot, linux or edk2). In each category can be any number of modules.

Each module has a name, which can be anything as long as it is unique (unique string across all modules of all types). In the example above there are 3 modules (`coreboot-example`, `linux-example`, `edk2-example`).

The configuration above can be simplified to this:
```
/
├── coreboot/
│   └── coreboot_example
├── linux/
│   └── linux_example
└── edk2/
    └── edk2_example
```

You can have multiple modules, as long as names are unique:
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
{{#include ../../../action/recipes/coreboot.go:CorebootOpts}}
```

`common` and `specific` are identical in function. There is no real difference between these two. They are split to simplify the code.

`depends` on the other hand allows you to specify dependency (or relation) between modules. For example your `coreboot` uses `edk2` as payload. So you can specify this dependency by listing name of the `edk2` module in `depends` of your `coreboot` module.

```json
{
  "coreboot": {
    "coreboot-example": {
      "depends": ["edk2-example"],
      "common": {
        ...
      },
      "specific": {
        ...
      }
    }
  },
  "edk2": {
    "edk2-example": {
      "depends": null,
      "common": {
        ...
      },
      "specific": {
        ...
      }
    }
  }
}
```

With such configuration, you can then run `firmware-action` recursively, and it will build all of the modules in proper order.
```
./firmware-action build --config=./my-config.json --target=coreboot-example --recursive
```
In this case `firmware-action` would build `edk2-example` and then `coreboot-example`.

```admonish tip
By changing inputs and outputs, you can then feed output of one module into inputs of another module.

This way you can build the entire firmware stack in single step.
```


## Common and Specific

To explain each and every entry in the configuration, here are snippets of the code with comments.

```admonish info
In the code below, the tag `json` (for example `:"sdk_url"`) specifies what the field is called in JSON file.

Tag `validate:"required"`, it means that the field is required and must not be empty. Empty required field will fail validation and terminate program with error.

Tag `validate:"dirpath"` means that field must contain a valid path to a directory. It is not necessary for the path or directory to exists, but must be a valid path. Be warned - that means that the string must end with `/`. For example `output-coreboot/`.

Tag `validate:"filepath"` means that the field must contain a valid path to a file. It is not necessary for the file to exist.

For more tails see [go-playground/validator](https://github.com/go-playground/validator) package.
```

### Common
```go
{{#include ../../../action/recipes/config.go:CommonOpts}}
```

### Specific / Coreboot
```go
{{#include ../../../action/recipes/coreboot.go:CorebootSpecific}}
```

### Specific / Linux
```go
{{#include ../../../action/recipes/linux.go:LinuxSpecific}}
```

### Specific / Edk2
```go
{{#include ../../../action/recipes/edk2.go:Edk2Specific}}
```

