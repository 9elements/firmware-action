# Docker containers

Docker is used to build the firmware stacks. To do this efficiently, purpose-specific docker containers are pre-build and are [published as packages](https://github.com/orgs/9elements/packages?repo_name=firmware-action) in GitHub repository.

However there was a problem with too many dockerfiles with practically identical content, just because of different version of software installed inside.

So to simplify this, we needed some master-configuration on top of our dockerfiles. Instead of making up some custom configuration solution, we just decided to use existing and defined [docker-compose](https://docs.docker.com/compose/) `yaml` config structure, with a custom parser (because there is no existing parser out there).

