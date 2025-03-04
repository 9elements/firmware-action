# Build Docker container on the fly

As already mentioned in [Configuration/common](config.md#common) section, `firmware-action` can build a Docker container on the fly when provided with `Dockerfile`.

The `sdk_url` field in configuration file accepts both URL and file-path. If file-path is provided, the container will be build and used (subsequent runs will no rebuild the container unless there changes were made to the `Dockerfile`).

The file-path can be a absolute or relative path to `Dockerfile` (or directory in which `Dockerfile` is) to build the image on the fly.

```admonish example title="Accepted values"
~~~
{{#include ../../../cmd/firmware-action/recipes/config.go:CommonOptsSdkURLExamples}}
~~~
```

```admonish warning
`file://` path cannot contain `..`
```

```admonish
Docker engine assumes single `Dockerfile` per directory, hence it requires path to the parent directory in which the `Dockerfile` resides (not to the file itself). For user-comfort, `firmware-action` accepts both path to parent directory and path to the file.

If the path contains the `Dockerfile` as last element, it will be removed before passed over to Docker engine.

Meaning that if user provides `file:///home/user/my-image/Dockerfile`, the Docker engine will receive `file:///home/user/my-image/`.
```
