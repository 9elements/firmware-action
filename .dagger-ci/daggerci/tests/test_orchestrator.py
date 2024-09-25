# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# pylint: disable=unused-import
# mypy: disable-error-code="import, no-untyped-def"

import datetime
import re
from unittest.mock import MagicMock, patch

import dagger
import pytest
from lib.docker_compose import DockerComposeValidate
from lib.orchestrator import ContainerMissingTestEnvVar


@pytest.mark.slow
def test__orchestrator__broken_compose(
    tmpdir, create_orchestrator, docker_compose_file_broken
):
    """
    Test with broken docker-compose file
    """
    with pytest.raises(DockerComposeValidate):
        create_orchestrator(
            dirpath=tmpdir, compose_file_content=docker_compose_file_broken
        )


@pytest.mark.slow
@pytest.mark.anyio
async def test__orchestrator__broken_dockerfile(
    tmpdir, create_orchestrator, dockerfile_broken
):
    """
    Test with broken dockerfile
    """
    my_orchestrator = create_orchestrator(
        dirpath=tmpdir, dockerfile_content=dockerfile_broken
    )
    result = await my_orchestrator.build_test_publish()
    assert "services" in result.results
    assert "coreboot_4.19" in result.results["services"]

    # because of multi-platform nature we have to be flexible
    build_found = False
    for key, _ in result.results["services"]["coreboot_4.19"].items():
        if re.match("build .*", key) and not re.match(".*_msg$", key):
            build_found = True
            assert result.results["services"]["coreboot_4.19"][key] is False
    assert build_found


@pytest.mark.slow
@pytest.mark.anyio
@pytest.mark.skip(reason="testing currently causes running out of disk space")
async def test__orchestrator__missing_env_var(tmpdir, create_orchestrator, dockerfile):
    """
    Try to execute test inside docker container, but there is no
        "VERIFICATION_TEST" defined
    """
    my_orchestrator = create_orchestrator(dirpath=tmpdir, dockerfile_content=dockerfile)
    with pytest.raises(ContainerMissingTestEnvVar):
        await my_orchestrator.build_test_publish()


@pytest.mark.slow
@pytest.mark.anyio
@pytest.mark.skip(reason="testing currently causes running out of disk space")
async def test__orchestrator__run_test_script_fail(
    tmpdir, create_orchestrator, dockerfile_dummy_tests_fail
):
    """
    Test container by running a script inside it,
        the script returns non-zero return code
    """
    my_orchestrator = create_orchestrator(
        dirpath=tmpdir, dockerfile_content=dockerfile_dummy_tests_fail
    )
    results = await my_orchestrator.build_test_publish()
    assert results.results == {
        "services": {
            "coreboot_4.19": {
                "build": True,
                "build_msg": None,
                "export": True,
                "export_msg": None,
                "test": False,
                "test_msg": None,
            }
        }
    }
    assert results.return_code == 1


@pytest.mark.slow
@pytest.mark.anyio
@pytest.mark.skip(reason="testing currently causes running out of disk space")
async def test__orchestrator__run_test_script_success(
    tmpdir, create_orchestrator, dockerfile_dummy_tests_success
):
    """
    Test container by running a script inside it
    """
    my_orchestrator = create_orchestrator(
        dirpath=tmpdir, dockerfile_content=dockerfile_dummy_tests_success
    )
    results = await my_orchestrator.build_test_publish()
    assert results.results == {
        "services": {
            "coreboot_4.19": {
                "build": True,
                "build_msg": None,
                "export": True,
                "export_msg": None,
                "test": True,
                "test_msg": None,
                "publish": False,
                "publish_msg": "skip",
            }
        }
    }
    assert results.return_code == 0


@pytest.mark.slow
@pytest.mark.anyio
@pytest.mark.skip(reason="testing currently causes running out of disk space")
async def test__orchestrator__multi_comprehensive_build(
    tmpdir,
    create_orchestrator,
    docker_compose_file_multi_comprehensive_build,
    dockerfile_dummy_tests_success,
):
    """
    Test the orchestrator with something similar to real-world use case
    """
    compose_yaml, expected_results = docker_compose_file_multi_comprehensive_build
    my_orchestrator = create_orchestrator(
        dirpath=tmpdir,
        compose_file_content=compose_yaml,
        dockerfile_content=dockerfile_dummy_tests_success,
    )
    results = await my_orchestrator.build_test_publish()
    assert results.results == expected_results


@pytest.mark.parametrize(
    "git_sha_short, git_sha_long, git_tag, git_branch, git_pull_request_number, expected",
    [
        (
            "sha12ab",
            "sha12ab0123456789abcdef0123456789abcdef0",
            None,
            "branch1",
            None,
            ["sha12ab", "branch1"],
        ),
        (
            "sha12ab",
            "sha12ab0123456789abcdef0123456789abcdef0",
            "v1.0",
            "master",
            None,
            ["sha12ab", "v1.0", "master", "latest"],
        ),
        (
            "sha12ab",
            "sha12ab0123456789abcdef0123456789abcdef0",
            "v1.1.0-3-11dd3b9",
            "branch1",
            None,
            ["sha12ab", "branch1"],
        ),
        (
            "sha12ab",
            "sha12ab0123456789abcdef0123456789abcdef0",
            "v1.1.0-3-11dd3b9",
            "branch1",
            "34",
            ["sha12ab", "branch1", "pull_request_34"],
        ),
        (
            "sha12ab",
            "sha12ab0123456789abcdef0123456789abcdef0",
            "v1.1.0-3-11dd3b9",
            "feature/new-stuff",
            "34",
            ["sha12ab", "feature_new_stuff", "pull_request_34"],
        ),
    ],
)
@pytest.mark.slow
def test__orchestrator__tags(
    tmpdir,
    create_orchestrator,
    monkeypatch,
    git_sha_short,
    git_sha_long,
    git_tag,
    git_branch,
    git_pull_request_number,
    expected,
):
    """
    Test tags
    """
    datetime_mock = MagicMock(wraps=datetime.datetime)
    FAKE_NOW = datetime.datetime(2023, 1, 1, 10, 20, 30)
    datetime_mock.now.return_value = FAKE_NOW
    monkeypatch.setattr(datetime, "datetime", datetime_mock)
    assert datetime.datetime.now() == FAKE_NOW

    with (
        patch("lib.git.git_describe", return_value=git_tag),
        patch(
            "lib.orchestrator.git_get_latest_commit_sha_short",
            return_value=git_sha_short,
        ),
        patch(
            "lib.orchestrator.git_get_latest_commit_sha_long", return_value=git_sha_long
        ),
        patch("lib.orchestrator.git_get_branch_name", return_value=git_branch),
        patch(
            "lib.orchestrator.github_get_pull_request_number",
            return_value=git_pull_request_number,
        ),
    ):
        my_orchestrator = create_orchestrator(dirpath=tmpdir)
    assert sorted(my_orchestrator.tags) == sorted(expected)
