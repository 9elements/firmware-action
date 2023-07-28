#!/usr/bin/python

# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments

import os
import textwrap
import pytest


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
