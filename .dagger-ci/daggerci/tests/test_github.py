# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import pytest
from lib.github import github_get_project_url, github_get_pull_request_number


@pytest.mark.parametrize(
    "env_vars, expected",
    [
        (
            #
            [
                ["GITHUB_SERVER_URL", "https://github.com"],
                ["GITHUB_REPOSITORY", "9elements/firmware-action"],
            ],
            "https://github.com/9elements/firmware-action",
        ),
        (
            [
                ["GITHUB_SERVER_URL", "https://git.com"],
                ["GITHUB_REPOSITORY", "whatever-repo"],
            ],
            "https://git.com/whatever-repo",
        ),
        (
            # Use fallback values
            [
                ["GITHUB_SERVER_URL", None],
                ["GITHUB_REPOSITORY", None],
            ],
            "https://github.com/9elements/firmware-action",
        ),
    ],
)
def test__github_get_project_url(monkeypatch, env_vars, expected):
    # Prepare environment
    for pair in env_vars:
        if pair[1] is None:
            monkeypatch.delenv(pair[0], raising=False)
        else:
            monkeypatch.setenv(pair[0], pair[1])
    # Test
    assert github_get_project_url() == expected


@pytest.mark.parametrize(
    "env_vars, expected",
    [
        (
            [["GITHUB_REF", "refs/heads/main"]],
            None,
        ),
        (
            [["GITHUB_REF", "refs/pull/34/merge"]],
            "34",
        ),
        (
            [["GITHUB_REF", "refs/tags/v1.0"]],
            None,
        ),
        (
            [["GITHUB_REF", None]],
            None,
        ),
    ],
)
def test__github_get_pull_request_number(monkeypatch, env_vars, expected):
    # Prepare environment
    for pair in env_vars:
        if pair[1] is None:
            monkeypatch.delenv(pair[0], raising=False)
        else:
            monkeypatch.setenv(pair[0], pair[1])
    # Test
    assert github_get_pull_request_number() == expected
