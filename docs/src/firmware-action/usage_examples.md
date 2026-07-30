# Examples

There is a separate directory with examples called [firmware-action-example](https://github.com/9elements/firmware-action-example/). It is a great source of information to get started.

In addition, the main repository also contains multiple examples (even though rather simple ones) defined in [.github/workflows/example.yml](https://github.com/9elements/firmware-action/blob/main/.github/workflows/example.yml). These are there to function as tests to verify the functionality, as such they are made with this specific task in mind. Please take that into account when going though them.


## Coreboot

> [!TIP]
> **Coreboot - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_coreboot}}
> ~~~

> [!TIP]
> **Coreboot - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__coreboot.json}}
> ~~~


## Linux Kernel

> [!TIP]
> **Linux Kernel - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_linux_kernel}}
> ~~~

> [!TIP]
> **Linux Kernel - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__linux.json}}
> ~~~


## Edk2

> [!TIP]
> **Edk2 - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_edk2}}
> ~~~

> [!TIP]
> **Edk2 - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__edk2.json}}
> ~~~


## Firmware Stitching

> [!TIP]
> **Firmware Stitching - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_stitch}}
> ~~~

> [!TIP]
> **Firmware Stitching - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__firmware_stitching.json}}
> ~~~


## u-root

> [!TIP]
> **u-root - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_uroot}}
> ~~~

> [!TIP]
> **u-root - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__uroot.json}}
> ~~~


## u-boot

> [!TIP]
> **u-boot - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_uboot}}
> ~~~

> [!TIP]
> **u-boot - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__uboot.json}}
> ~~~


## Universal

> [!TIP]
> **Universal - GitHub CI**
>
> ~~~yaml
> {{#include ../../../.github/workflows/example.yml:example_build_universal}}
> ~~~

> [!TIP]
> **Universal - Configuration file**
>
> ~~~json
> {{#include ../../../tests/example_config__universal.json}}
> ~~~

