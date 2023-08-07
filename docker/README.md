# Docker related files

Collection of dockerfiles and open source tools to
build various firmware images.


## Docker-compose

The compose file is not used at the moment by actual docker-compose, it is manually parsed and fed to dagger. In addition, at the time of writing, dagger does not support docker-compose.

The core problem was too many dockerfiles with practically identical content. And instead ot making up some custom configuration solution, we just decided to use existing and defined docker-compose yaml config structure.

The yaml file is parsed by naive and simplistic parser which supports bare-minimum features just to get parameters needed by dagger to build all docker containers.

Example of yaml config file:
```yaml
services:
  coreboot_4.19:
    build:
      context: coreboot
      args:
        - COREBOOT_VERSION=4.19
  coreboot_4.20:
    build:
      context: coreboot
      args:
        - COREBOOT_VERSION=4.20
```


## Environment variables

There are environment variables `VERIFICATION_TEST` and possibly more `VERIFICATION_TEST_*`, these are for when the docker containers are built to test the functionality.

`VERIFICATION_TEST` should contain a path to test script, `VERIFICATION_TEST_*` variables are used in said script.

