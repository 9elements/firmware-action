# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import re
import pytest
from unittest.mock import patch

from lib.git import (
    git_get_latest_commit_sha_long,
    git_get_latest_commit_sha_short,
    git_describe,
    git_get_tag,
)


def test__git_get_latest_commit_sha_long():
    assert re.match(r"^[a-z\d]{40}$", git_get_latest_commit_sha_long())


def test__git_get_latest_commit_sha_short():
    assert re.match(r"^[a-z\d]{7}$", git_get_latest_commit_sha_short())


@pytest.mark.parametrize(
    "describe, expected",
    [
        ("1.0", "1.0"),
        ("v1.1", "v1.1"),
        ("v1.1.0", "v1.1.0"),
        ("1.1.0", "1.1.0"),
        ("1.1.0-3-11dd3b9", None),
        ("v1.1.0-3-11dd3b9", None),
        ("lol", None),
        (None, None),
    ],
)
def test__git_get_tag(describe, expected):
    with patch("lib.git.git_describe", wraps=git_describe) as mock_git_describe:
        mock_git_describe.return_value = describe
        assert git_get_tag() == expected
