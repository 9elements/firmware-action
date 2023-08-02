#!/usr/bin/python

# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments

import os
import textwrap
import pytest

from lib.filesystem import mkdir


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
    '''
    Needed for anyio
      https://anyio.readthedocs.io/en/stable/testing.html#testing-with-anyio
    '''
    return 'asyncio'


# ===========================
#
#  Common fixtures
#
# ===========================

@pytest.fixture
def create_file():
    '''
    Create text file at "path" with "content" as its content.
    '''
    def _create_file(path: str, content: str):
        rootdir = os.path.dirname(path)
        if not os.path.isdir(rootdir):
            mkdir(rootdir)
        with open(path, 'w', encoding='utf-8') as myfile:
            myfile.write(content)
    return _create_file


# ===========================
#
#  Dockerfile fixtures
#
# ===========================

@pytest.fixture
def dockerfile(create_file):
    '''
    Generic Dockerfile content
    '''
    return textwrap.dedent("""\
        FROM ubuntu:22.04 AS base
        ARG TARGETARCH=amd64
        ARG COREBOOT_VERSION=4.19
        RUN apt-get update && \\
            apt-get install -y --no-install-recommends \\
                bc nano git \\
            && \\
            rm -rf /var/lib/apt/lists/*\
            """)


@pytest.fixture
def dockerfile_dummy_tests(create_file):
    '''
    Dockerfile content specifically for executing tests inside docker
    '''
    return textwrap.dedent("""\
        FROM ubuntu:22.04 AS base
        ARG TARGETARCH=amd64
        ARG CONTEXT=dummy
        ARG VARIANT=success
        ENV VERIFICATION_TEST=./tests/test_${CONTEXT}_${VARIANT}.sh
        RUN apt-get update && \\
            apt-get install -y --no-install-recommends \\
                bc nano git \\
            && \\
            rm -rf /var/lib/apt/lists/*\
            """)


@pytest.fixture
def dockerfile_broken(create_file):
    '''
    Dockerfile content which should fail to build
    '''
    return textwrap.dedent("""\
        FROM ubuntu:22.04 AS base
        RUN false\
            """)


# ===========================
#
#  Docker Compose fixtures
#
# ===========================

@ pytest.fixture
def docker_compose_file(create_file):
    '''
    Generic Docker compose
    '''
    return textwrap.dedent("""\
        services:
          coreboot_4.19:
            build:
              context: coreboot""")


@ pytest.fixture
def docker_compose_file_broken(create_file):
    '''
    Docker compose which should fail syntax falidation
    '''
    return textwrap.dedent("""\
        services:
          coreboot_4.19:
            asdfgh context coreboot""")


@ pytest.fixture
def docker_compose_file_complex(create_file):
    # TODO
    return textwrap.dedent("""\
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
        """)
