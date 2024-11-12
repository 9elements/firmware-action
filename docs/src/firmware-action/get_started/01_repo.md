# Start a new git repository

Start a new repository in GitHub and then clone it.


## Add coreboot as git submodule

```admonish tip
To add git submodule, run:
~~~
git submodule add <repo> <path>
~~~
```

Add [coreboot repository](https://review.coreboot.org/admin/repos/coreboot,general) as a submodule:
```bash
git submodule add --depth=1 "https://review.coreboot.org/coreboot" coreboot
```

In this example we will work with coreboot `4.19` release (it is a bit older release from January 2023, but should suffice for demonstration)
```bash
( cd coreboot; git fetch origin tag "4.19"; git checkout "4.19" )
```

Recursively initialize submodules.

```bash
git submodule update --init --recursive
```

```admonish warning
Recursively initializing all submodules in coreboot will take a minute or two.
```

