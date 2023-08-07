# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import os
import textwrap

import pytest
from lib.docker_compose import DockerCompose
from lib.filesystem import mkdir
from lib.orchestrator import Orchestrator

# ===========================
#  Add SLOW marker
# ===========================

# Tests marked with "@pytest.mark.slow" are automatically skipped
# To tun slow tests, call pytest with additional "--runslow" argument
# Docs: https://docs.pytest.org/en/latest/example/simple.html#control-skipping-of-tests-according-to-command-line-option


def pytest_addoption(parser):
    parser.addoption(
        "--runslow", action="store_true", default=False, help="run slow tests"
    )


def pytest_configure(config):
    config.addinivalue_line("markers", "slow: mark test as slow to run")


def pytest_collection_modifyitems(config, items):
    if config.getoption("--runslow"):
        # --runslow given in cli: do not skip slow tests
        return
    skip_slow = pytest.mark.skip(reason="need --runslow option to run")
    for item in items:
        if "slow" in item.keywords:
            item.add_marker(skip_slow)


# ===========================
#  Allow testing
# ===========================


@pytest.fixture
def anyio_backend():
    """
    Needed for anyio
      https://anyio.readthedocs.io/en/stable/testing.html#testing-with-anyio
    """
    return "asyncio"


# ===========================
#
#  Common fixtures
#
# ===========================


@pytest.fixture
def create_file():
    """
    Create text file at "path" with "content" as its content.
    """

    def _create_file(path: str, content: str):
        rootdir = os.path.dirname(path)
        if not os.path.isdir(rootdir):
            mkdir(rootdir)
        with open(path, "w", encoding="utf-8") as myfile:
            myfile.write(content)

    return _create_file


# ===========================
#
#  Dockerfile fixtures
#
# ===========================


@pytest.fixture
def dockerfile():
    """
    Generic Dockerfile content
    """
    return textwrap.dedent(
        """\
        FROM ubuntu:22.04 AS base
        ARG TARGETARCH=amd64
        ARG COREBOOT_VERSION=4.19
        RUN apt-get update && \\
            apt-get install -y --no-install-recommends \\
                bc nano git \\
            && \\
            rm -rf /var/lib/apt/lists/*\
            """
    )


@pytest.fixture
def dockerfile_dummy_tests_success():
    """
    Dockerfile content specifically for executing tests inside docker
    """
    return textwrap.dedent(
        """\
        FROM ubuntu:22.04 AS base
        ARG TARGETARCH=amd64
        ARG CONTEXT=dummy
        ARG VARIANT=success
        ENV VERIFICATION_TEST=./tests/test_${CONTEXT}_${VARIANT}.sh
        RUN echo 'hello world'\
        """
    )


@pytest.fixture
def dockerfile_dummy_tests_fail():
    """
    Dockerfile content specifically for executing tests inside docker
    """
    return textwrap.dedent(
        """\
        FROM ubuntu:22.04 AS base
        ARG TARGETARCH=amd64
        ARG CONTEXT=dummy
        ARG VARIANT=fail
        ENV VERIFICATION_TEST=./tests/test_${CONTEXT}_${VARIANT}.sh
        RUN echo 'hello world'\
        """
    )


@pytest.fixture
def dockerfile_broken():
    """
    Dockerfile content which should fail to build
    """
    return textwrap.dedent(
        """\
        FROM ubuntu:22.04 AS base
        RUN false\
            """
    )


# ===========================
#
#  Docker Compose fixtures
#
# ===========================


@pytest.fixture
def docker_compose_file():
    """
    Generic Docker compose
    """
    return textwrap.dedent(
        """\
        services:
          coreboot_4.19:
            build:
              context: coreboot"""
    )


@pytest.fixture
def docker_compose_file_broken():
    """
    Docker compose which should fail syntax validation
    """
    return textwrap.dedent(
        """\
        services:
          coreboot_4.19:
            asdfgh context coreboot"""
    )


@pytest.fixture
def docker_compose_file_complex():
    # TODO
    return textwrap.dedent(
        """\
        services:
          coreboot_4.19:
            build:
              context: coreboot
              args:
                - COREBOOT_VERSION=4.19
          coreboot_4.20:
            build:
              args:
                - COREBOOT_VERSION=4.20
          edk2:
            build:
              context: edk2
          meh:
            build:
              args:
                - more=meh
          meh2:
            image: ubuntu\
        """
    )


@pytest.fixture
def docker_compose_file_multi_comprehensive_build():
    return [
        textwrap.dedent(
            """\
        services:
          dummy_1:
            build:
              context: dummy
              args:
                - VARIANT=success
          dummy_2:
            build:
              context: dummy
              args:
                - VARIANT=success
          dummy_3:
            build:
              context: dummy
              args:
                - VARIANT=fail
        """
        ),
        {
            "services": {
                "dummy_1": {
                    "build": True,
                    "build_msg": None,
                    "export": True,
                    "export_msg": None,
                    "test": True,
                    "test_msg": None,
                    "publish": False,
                    "publish_msg": "skip",
                },
                "dummy_2": {
                    "build": True,
                    "build_msg": None,
                    "export": True,
                    "export_msg": None,
                    "test": True,
                    "test_msg": None,
                    "publish": False,
                    "publish_msg": "skip",
                },
                "dummy_3": {
                    "build": True,
                    "build_msg": None,
                    "export": True,
                    "export_msg": None,
                    "test": False,
                    "test_msg": None,
                },
            }
        },
    ]


# ===========================
#
#  Misc fixtures
#
# ===========================


@pytest.fixture
def create_orchestrator(create_file, docker_compose_file, dockerfile):
    def _create_orchestrator(
        dirpath: str,
        compose_file_content: str | None = None,
        dockerfile_content: str | None = None,
    ):
        # Create docker compose
        docker_compose_file_path = os.path.join(dirpath, "compose.yaml")
        if compose_file_content is None:
            compose_file_content = docker_compose_file
        create_file(path=docker_compose_file_path, content=compose_file_content)

        # Create dockerfiles according to docker compose
        my_dockercompose = DockerCompose(path=docker_compose_file_path)
        for df in my_dockercompose.get_dockerfiles():
            dockerfile_path = os.path.join(
                dirpath, my_dockercompose.get_dockerfile_context(df), "Dockerfile"
            )
            if dockerfile_content is None:
                dockerfile_content = dockerfile
            if not os.path.isfile(dockerfile_path):
                create_file(path=dockerfile_path, content=dockerfile_content)

        return Orchestrator(docker_compose_path=docker_compose_file_path)

    return _create_orchestrator
