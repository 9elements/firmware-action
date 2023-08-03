# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import os
import pytest

from lib.env import get_env_var_value


@pytest.mark.parametrize(
    "list_of_maybe_existing_env_vars, fallback, presets, expected",
    [
        # No env vars set, should return fallfack
        (['NONEXISTING_A', 'NONEXISTING_B'], '/tmp', [], '/tmp'),
        # First env var set, should return that
        (['NONEXISTING_A', 'NONEXISTING_B'], '/tmp', [['NONEXISTING_A', '/home']], '/home'),
        # Two env vars set, should return first one
        (['NONEXISTING_A', 'NONEXISTING_B'], '/tmp',
         [['NONEXISTING_A', '/home'], ['NONEXISTING_B', '/usr']], '/home'),
    ])
def test__get_env_var_value(monkeypatch, list_of_maybe_existing_env_vars, fallback, presets, expected):
    # Prepare environment
    for pair in presets:
        monkeypatch.setenv(pair[0], pair[1])
    # Test
    assert get_env_var_value(
        list_of_maybe_existing_env_vars=list_of_maybe_existing_env_vars, fallback=fallback) == expected
