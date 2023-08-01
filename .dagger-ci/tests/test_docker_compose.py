# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments

import os
import dagger
import pytest
from contextlib import nullcontext as does_not_raise

from lib.docker_compose import *


@pytest.mark.parametrize(
    "heap, needle, expected, expectation",
    [
        (['a', 'b'], 'a',  'a', does_not_raise()),
        (['a', 'b'], 'b',  'b', does_not_raise()),
        (['a', 'b'], None, 'a', does_not_raise()),
        (['a', 'b'], 'c',  'a', pytest.raises(ValueError)),
    ])
def test__select(heap, needle, expected, expectation):
    with expectation:
        assert select(heap=heap, needle=needle) == expected


def test__docker_compose_broken(tmpdir, create_docker_compose_file_broken):
    compose_file = os.path.join(tmpdir, 'compose.yaml')
    create_docker_compose_file_broken(path=compose_file)
    with pytest.raises(DockerComposeValidate):
        mydockercompose = DockerCompose(path=compose_file)


def test__docker_compose(tmpdir, create_docker_compose_file_complex):
    compose_file = os.path.join(tmpdir, 'compose.yaml')
    create_docker_compose_file_complex(path=compose_file)
    mydockercompose = DockerCompose(path=compose_file)

    assert mydockercompose.get_top_elements() == ['services']

    assert mydockercompose.__select_top_element__() == 'services'
    assert mydockercompose.__select_top_element__('services') == 'services'
    with pytest.raises(ValueError):
        mydockercompose.__select_top_element__('')

    assert mydockercompose.get_dockerfiles() == ['coreboot_4.19', 'coreboot_4.20', 'edk2', 'meh']
    assert mydockercompose.get_dockerfiles(
        top_element='services') == ['coreboot_4.19', 'coreboot_4.20', 'edk2', 'meh']

    assert mydockercompose.__select_dockerfile__() == 'coreboot_4.19'
    assert mydockercompose.__select_dockerfile__('coreboot_4.19') == 'coreboot_4.19'
    assert mydockercompose.__select_dockerfile__(
        'coreboot_4.19', 'services') == 'coreboot_4.19'
    with pytest.raises(ValueError):
        mydockercompose.__select_dockerfile__('')
    with pytest.raises(ValueError):
        mydockercompose.__select_dockerfile__('', '')

    assert mydockercompose.get_dockerfile_context() == 'coreboot'
    assert mydockercompose.get_dockerfile_context('coreboot_4.19') == 'coreboot'
    assert mydockercompose.get_dockerfile_context('coreboot_4.19', 'services') == 'coreboot'
    assert mydockercompose.get_dockerfile_context('coreboot_4.20') == None
    assert mydockercompose.get_dockerfile_context('meh') == None

    assert mydockercompose.get_dockerfile_args(
    ) == [dagger.api.gen.BuildArg('COREBOOT_VERSION', '4.19')]
    assert mydockercompose.get_dockerfile_args('coreboot_4.19') == [
        dagger.api.gen.BuildArg('COREBOOT_VERSION', '4.19')]
    assert mydockercompose.get_dockerfile_args(
        'coreboot_4.19', 'services') == [dagger.api.gen.BuildArg('COREBOOT_VERSION', '4.19')]
    assert mydockercompose.get_dockerfile_args('edk2') == []
