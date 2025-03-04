# Change detection

`firmware-action` has basic detection of changes.

Related temporary files to aid in change detection are stored in `.firmware-action/` directory, which is always created in current working directory.

```admonish note
It is save to delete the `.firmware-action/` directory, but keep in mind that change detection depends on its existence.

If `.firmware-action/` directory is deleted, it will be re-created on next `firmware-action` execution.

It might be advantageous to also include this directory in caches / artifacts when in CI. Preserving these files might reduce run time in the CI.
```


## Sources modification time

When building a module, `firmware-action` checks recursively all sources for a given module. For all module types, list of sources include repository and all input files.

```admonish example collapsible=true title="Code snippet: Common sources"
~~~golang
{{#include ../../../cmd/firmware-action/recipes/config.go:CommonOptsGetSources}}
~~~
```

Each module type has then additional source files. For example `coreboot`, where list of sources also includes `defconfig` and all the blobs.
```admonish example collapsible=true title="Code snippet: Additional coreboot sources"
~~~golang
{{#include ../../../cmd/firmware-action/recipes/coreboot.go:CorebootOptsGetSources}}
~~~
```

When a module is successfully built, a file containing time stamp is saved to `.firmware-action/timestamps/` directory.

On next run, this file (if exists) is loaded with time stamp of last successful run. Then all sources are recursively checked for any file that was modified since the last successful run. If no file was modified since the loaded time stamp, module is considered up-to-date and build is skipped. If any of the files has newer modified time, module is re-built.


## Configuration file changes

`firmware-action` can also detect changes in the configuration file. For each module, on each successful build, it stores a copy of the configuration in `.firmware-action/configs/` directory.

On next run, current configuration is compared to configuration of last successful build, and if the configuration for the specific module differs, module is re-built.


## Git commit hash changes

`firmware-action` can detect changes based on git commit hashes. For each module, on each successful build, it stores the git commit hash of the module's repository path (`repo_path`) in `.firmware-action/git-hash/` directory.

On next run, the current git commit hash of the module's repository is compared to the stored hash from the last successful build. If the hashes differ, indicating that the module's repository has been changed, the module is re-built.
