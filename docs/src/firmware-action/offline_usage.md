# Offline usage


## Occasionally offline machines

This approach is useful for example when people travel. They run `firmware-action` beforehand to pull the necessary containers and then use this method to work offline.

`firmware-action` under the hood uses [dagger](https://docs.dagger.io/) / docker. As such, the configuration contains entry `sdk_url` which specifies the docker image / container to use.

```admonish example
~~~json
"sdk_url": "ghcr.io/9elements/firmware-action/coreboot_4.19:main"
~~~
```

However, in this configuration `firmware-action` (or rather `dagger`) will always connect to the internet and download the manifest to see if a new container needs to be downloaded. This applies to all tags (`main`, `latest`, `v0.8.0`, and so on).

If you need to use `firmware-action` offline, you have to first acquire the container. Either by running `firmware-action` at least once online, or by other means provided by docker.

Then you need to change the `firmware-action` configuration to include the image reference (digest hash).

```admonish example
~~~json
"sdk_url": "http://ghcr.io/9elements/firmware-action/coreboot_4.19:main@sha256:25b4f859e26f84a276fe0c4395a4f0c713f5b564679fbff51a621903712a695b"
~~~
```

Digest hash can be found in the container hub. For `firmware-action` containers see [GitHub](https://github.com/orgs/9elements/packages?repo_name=firmware-action).

It will also be displayed every time `firmware-action` is executed as `INFO` message near the start:
```
[INFO   ] Container information
    - time: 2024-12-01T12:09:43.62620859+01:00
    - Image reference: ghcr.io/9elements/firmware-action/coreboot_4.19:main@sha256:25b4f859e26f84a276fe0c4395a4f0c713f5b564679fbff51a621903712a695b
    - origin of this message: container.Setup
```

Simply copy-paste the digest (or image reference) into your configuration file and `firmware-action` will not connect to the internet to fetch a container if one matching is already present.


## Always offline machines

Besides running a offline docker registry, there is also a option to use tarballs.

`firmware-action` can import a tarfile and use it.

```admonish example
~~~json
"sdk_url": "file:///home/user/my-image/ubuntu-latest.tar"
~~~
```
