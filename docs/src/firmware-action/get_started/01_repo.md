# Start a new git repository

Start a new repository in GitHub and the clone it.


## Add coreboot as submodule

~~~admonish tip
Add git submodule with:
```
git submodule add <repo> <path>
```
~~~

Add [coreboot repository](https://review.coreboot.org/admin/repos/coreboot,general) as a submodule:
```bash
git submodule add --depth=1 "https://review.coreboot.org/coreboot" coreboot
```

```bash
git submodule update --init
```

Optionally checkout a release tag, for example `4.19` (it is a bit older release from January 2023, but should suffice for demonstration)
```bash
( cd coreboot; git fetch origin tag "4.19"; git checkout "4.19" )
```

Recursively initialize submodules.

```bash
git submodule update --init --recursive
```

~~~admonish warning
Recursively initializing all submodules in coreboot will take a moment.
~~~

