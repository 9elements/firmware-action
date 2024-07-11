# How to generate defconfig

## coreboot

Make your changes:
```bash
$ make menuconfig
```

Generate `defconfig` file:
```bash
$ make savedefconfig
```


## Linux

```bash
$ make tinyconfig
$ make savedefconfig
```

