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
    def _create_file(path: str, content: str):
        rootdir = os.path.dirname(path)
        if not os.path.isdir(rootdir):
            mkdir(rootdir)
        with open(path, 'w', encoding='utf-8') as myfile:
            myfile.write(content)
    return _create_file


@pytest.fixture
def create_dockerfile(create_file):
    def _create_dockerfile(path: str):
        create_file(path=path, content=textwrap.dedent("""\
                FROM ubuntu:22.04 AS base
                ARG TARGETARCH=amd64
                ARG COREBOOT_VERSION=4.19
                RUN apt-get update && \\
                    apt-get install -y --no-install-recommends \\
                        bc nano git \\
                    && \\
                    rm -rf /var/lib/apt/lists/*\
                    """))
    return _create_dockerfile


@pytest.fixture
def create_dockerfile_broken(create_file):
    def _create_dockerfile_broken(path: str):
        create_file(path=path, content=textwrap.dedent("""\
                FROM ubuntu:22.04 AS base
                RUN false\
                    """))
    return _create_dockerfile_broken


@pytest.fixture
def create_docker_compose_file(create_file):
    def _create_docker_compose_file(path: str):
        create_file(path=path, content=textwrap.dedent("""\
                services:
                  coreboot_4.19:
                    build:
                      context: coreboot"""))
    return _create_docker_compose_file


@pytest.fixture
def create_docker_compose_file_complex(create_file):
    def _create_docker_compose_file_complex(path: str):
        create_file(path=path, content=textwrap.dedent("""\
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
                        - more=meh\
                """))
    return _create_docker_compose_file_complex
