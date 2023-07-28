# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments

import os
import pytest

from lib.filesystem import *


def test__mkdir(tmpdir):
    path = os.path.join(tmpdir, 'whatever')

    # directory does not exist
    assert not os.path.isdir(path)
    mkdir(path)
    assert os.path.isdir(path)

    # directory does exist
    mkdir(path)
    assert os.path.isdir(path)
