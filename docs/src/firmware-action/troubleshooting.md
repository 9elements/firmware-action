# Troubleshooting common problems

Many `firmware-action` errors and warnings come with suggestion on how to fix them.

Other than that, here are some common problems and solutions.


```admonish tip
The first thing when troubleshooting is to look through the output for errors and warnings. Many of these messages come with a `suggestion` with instructions on possible solutions.

For example a warning message:
~~~
[WARN   ] Git submodule seems to be uninitialized
    - time: 2024-12-02T12:42:33.31416978+01:00
    - suggestion: run 'git submodule update --depth 0 --init --recursive --checkout'
    - offending_submodule: coreboot-linuxboot-example/linux
    - origin of this message: main.run
~~~
```


## Missing submodules / missing files

The problem can manifest in multiple way, most commonly with error messages of missing files.
```
make: *** BaseTools: No such file or directory.  Stop.
```

Solution is to get all git submodules.
```
git submodule update --depth 1 --init --recursive --checkout
```


## Coreboot blob not found

Blobs are copied into container separately from `input_files` and `input_dirs`, the path should point to files on your host.


## Dagger problems

To troubleshoot dagger, please see [dagger documentation](https://docs.dagger.io/troubleshooting).
