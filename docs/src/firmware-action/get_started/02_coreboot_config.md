# Create a coreboot configuration file

Now we need to create a configuration file for coreboot.

Either follow a [coreboot guide](https://doc.coreboot.org/tutorial/part1.html#step-5-configure-the-build) to get a bare-bones-basic configuration, or just copy-paste this text into `seabios_defconfig` file.

~~~admonish example title="seabios_defconfig"
```properties
{{#include ../../firmware-action-example/coreboot-example/coreboot_seabios_defconfig}}
```
~~~

