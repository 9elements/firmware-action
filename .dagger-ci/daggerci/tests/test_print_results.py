# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import pytest
from lib.results import Results


@pytest.mark.parametrize(
    "results, expected",
    [
        ([["services", "coreboot", "build", True]], 0),
        ([["services", "coreboot", "build", False]], 1),
        (
            [
                ["services", "coreboot", "build", True],
                ["services", "coreboot", "export", True],
                ["services", "coreboot", "test", True],
                ["services", "coreboot", "publish", False, "skip"],
            ],
            0,
        ),
        (
            [
                ["services", "coreboot", "build", True],
                ["services", "coreboot", "export", True],
                ["services", "coreboot", "test", True],
                ["services", "coreboot", "publish", False, "skip"],
                ["services", "edk2", "build", False],
            ],
            1,
        ),
    ],
)
def test__results(results, expected):
    my_results = Results()
    for i in results:
        my_results.add(*i)
    my_results.print()
    assert my_results.return_code == expected
