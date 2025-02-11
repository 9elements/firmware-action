# Docker-compose

The compose file is not used at the moment by actual [docker-compose](https://docs.docker.com/compose/), it is manually parsed and then fed to dagger.

```admonish info
dagger does not support docker-compose at the time of writing.

There is also no existing docker-compose parser out there that we could use as off-the-shelf solution.
```

The custom parser implements only a limited feature-set out of the [compose-file spec](https://docs.docker.com/compose/compose-file/), just the bare-minimum needed to build the containers:
- [services](https://docs.docker.com/compose/compose-file/05-services/)
    - ~~attach~~
    - [build](https://docs.docker.com/compose/compose-file/build/)
        - [context](https://docs.docker.com/compose/compose-file/build/#context)
        - ~~dockerfile~~
        - ~~dockerfile_inline~~
        - [args](https://docs.docker.com/compose/compose-file/build/#args)
        - ~~ssh~~
        - ...
    - ~~blkio_config~~
    - ...
- ~~networks~~
- ~~volumes~~
- ~~configs~~
- ~~secrets~~

This way, we can have a single parametric `Dockerfile` for each item (coreboot, linux, edk2, ...) and introduce variation with scalable and maintainable single-file config.

```admonish example
Example of `compose.yaml` file to build 2 versions of `coreboot` docker image:
~~~yaml
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
~~~
```



## Multi-stage builds

We use [multi-stage builds](https://docs.docker.com/build/building/multi-stage/) to minimize the final container / image.


## Environment variables

In the dockerfiles, we heavily rely on use of [environment variables](https://docs.docker.com/engine/reference/builder/#env) and [arguments](https://docs.docker.com/engine/reference/builder/#arg).

This allows for the parametric nature.


## Testing

Containers are also tested to verify that they were build successfully.

The tests are rather simple, consisting solely from [happy path](https://en.wikipedia.org/wiki/Happy_path) tests. This might change in the future.

Test is done by executing a shell script which builds firmware in some hello-world configuration example. Nothing too fancy.

The path to said shell script is stored in environment variable `VERIFICATION_TEST`.

```admonish example collapsible=true title="Example of coreboot test"
~~~bash
{{#include ../../../tests/test_coreboot.sh}}
~~~
```

In addition, there might be `VERIFICATION_TEST_*` variables. These are used inside the test script and are rather use-case specific, however often used to store which version of firmware is being tested.


## Adding new container

- (optional) Add new `Dockerfile` into `docker` directory
- Add new entry in `docker/compose.yaml`
- Add new entry into strategy matrix in `.github/workflows/docker-build-and-test.yml`
- (optional) Add new strategy matrix in `.github/workflows/example.yml` examples
    - this requires adding new configuration file in `tests` directory
- Add entry into a list of containers in `README.md`


## Discontinuing container

- Update entry in list of containers in `README.md`
- Add new regex entry into `Setup()` function in `cmd/firmware-action/container/container.go` to warn about discontinued containers
