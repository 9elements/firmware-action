# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

from contextlib import nullcontext as does_not_raise

import pytest
from lib.cli import cli


@pytest.mark.parametrize(
    "args, expectation",
    [
        (["--help"], pytest.raises(SystemExit)),
        ([], does_not_raise()),
        (["-d"], pytest.raises(SystemExit)),
    ],
)
def test__cli__smoke_test(args, expectation):
    with expectation:
        cli(args=args)


@pytest.mark.parametrize(
    "args, expected",
    [
        (["-c"], True),
        ([], False),
    ],
)
def test__cli__concurrent(args, expected):
    arguments = cli(args=args)
    assert arguments.concurrent == expected


@pytest.mark.parametrize(
    "args, expected",
    [
        (["-v"], True),
        ([], False),
    ],
)
def test__cli__verbose(args, expected):
    arguments = cli(args=args)
    assert arguments.verbose == expected


@pytest.mark.parametrize(
    "args, expected",
    [
        (["-p"], True),
        ([], False),
    ],
)
def test__cli__publish(args, expected):
    arguments = cli(args=args)
    assert arguments.publish == expected


@pytest.mark.parametrize(
    "args, expected",
    [
        (["-d", "hello"], ["hello"]),
        (["-d", "hello", "world"], ["hello", "world"]),
    ],
)
def test__cli__dockerfile(args, expected):
    arguments = cli(args=args)
    assert arguments.dockerfile == expected
