# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import os
from contextlib import nullcontext as does_not_raise

import dagger
import pytest
from lib.docker_compose import (
    DockerCompose,
    DockerComposeMissingElement,
    DockerComposeValidate,
    select,
)


@pytest.mark.parametrize(
    "heap, needle, expected, expectation",
    [
        (["a", "b"], "a", "a", does_not_raise()),
        (["a", "b"], "b", "b", does_not_raise()),
        (["a", "b"], None, "a", does_not_raise()),
        (["a", "b"], "c", "a", pytest.raises(ValueError)),
    ],
)
def test__select(heap, needle, expected, expectation):
    with expectation:
        assert select(heap=heap, needle=needle) == expected


def test__docker_compose_broken(tmpdir, create_file, docker_compose_file_broken):
    compose_file = os.path.join(tmpdir, "compose.yaml")
    create_file(path=compose_file, content=docker_compose_file_broken)
    with pytest.raises(DockerComposeValidate):
        DockerCompose(path=compose_file)


def test__docker_compose(tmpdir, create_file, docker_compose_file_complex):
    compose_file = os.path.join(tmpdir, "compose.yaml")
    create_file(path=compose_file, content=docker_compose_file_complex)
    my_docker_compose = DockerCompose(path=compose_file)

    assert my_docker_compose.get_top_elements() == ["services"]

    assert my_docker_compose.__select_top_element__() == "services"
    assert my_docker_compose.__select_top_element__("services") == "services"
    with pytest.raises(DockerComposeMissingElement):
        my_docker_compose.__select_top_element__("")

    assert my_docker_compose.get_dockerfiles() == [
        "coreboot_4.19",
        "coreboot_4.20",
        "edk2",
        "meh",
        "meh2",
    ]
    assert my_docker_compose.get_dockerfiles(top_element="services") == [
        "coreboot_4.19",
        "coreboot_4.20",
        "edk2",
        "meh",
        "meh2",
    ]

    assert my_docker_compose.__select_dockerfile__() == "coreboot_4.19"
    assert my_docker_compose.__select_dockerfile__("coreboot_4.19") == "coreboot_4.19"
    assert (
        my_docker_compose.__select_dockerfile__("coreboot_4.19", "services")
        == "coreboot_4.19"
    )
    with pytest.raises(DockerComposeMissingElement):
        my_docker_compose.__select_dockerfile__("")
    with pytest.raises(DockerComposeMissingElement):
        my_docker_compose.__select_dockerfile__("", "")

    assert my_docker_compose.get_dockerfile_context() == "coreboot"
    assert my_docker_compose.get_dockerfile_context("coreboot_4.19") == "coreboot"
    assert (
        my_docker_compose.get_dockerfile_context("coreboot_4.19", "services")
        == "coreboot"
    )
    assert my_docker_compose.get_dockerfile_context("coreboot_4.20") is None
    assert my_docker_compose.get_dockerfile_context("meh") is None
    assert my_docker_compose.get_dockerfile_context("meh2") is None

    assert my_docker_compose.get_dockerfile_args() == [
        dagger.api.gen.BuildArg("COREBOOT_VERSION", "4.19")
    ]
    assert my_docker_compose.get_dockerfile_args("coreboot_4.19") == [
        dagger.api.gen.BuildArg("COREBOOT_VERSION", "4.19")
    ]
    assert my_docker_compose.get_dockerfile_args("coreboot_4.19", "services") == [
        dagger.api.gen.BuildArg("COREBOOT_VERSION", "4.19")
    ]
    assert my_docker_compose.get_dockerfile_args("edk2") == []
