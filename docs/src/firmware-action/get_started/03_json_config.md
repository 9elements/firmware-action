# Create a JSON configuration file

This configuration file is for firmware-action, so that it knows what to do and where to find things. Let's call it `firmware-action.json`.

~~~admonish example title="firmware-action.json"
```json
{{#include ../../firmware-action-example/coreboot-example.json}}
```
~~~

~~~admonish info
Field `repo_path` is pointing to the location of our coreboot submodule which we added in previous step [Repository](01_repo.md).
~~~

~~~admonish info
Field `defconfig_path` is pointing to the location of coreboot's configuration file which we created in previous step [coreboot configuration](02_coreboot_config.md).
~~~

~~~admonish info
Firmware action can be used to compile other firmware too, and even combine multiple firmware projects (to a certain degree).

For this reason the JSON configuration file is divided into categories (`coreboot`, `edk2`, etc). Each category can contain multiple entries.

Entries can depend on each other, which allows you to combine them - you can have for example `coreboot` firmware with `edk2` payload.
~~~

