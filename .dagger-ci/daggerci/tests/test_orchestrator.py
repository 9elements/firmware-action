# pylint: disable=missing-function-docstring
# pylint: disable=missing-module-docstring
# pylint: disable=too-many-arguments
# mypy: disable-error-code="import, no-untyped-def"

import dagger
import pytest
from lib.docker_compose import DockerComposeValidate
from lib.orchestrator import ContainerMissingTestEnvVar


@pytest.mark.slow
@pytest.mark.anyio
async def test__orchestrator__broken_compose(
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
    assert result.results == {
        "services": {
            "coreboot_4.19": {
                "build": False,
                "build_msg": 'failed to solve: process "/bin/sh -c false" did not complete successfully: exit code: 1',
            }
        }
    }


@pytest.mark.slow
@pytest.mark.anyio
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
async def test__orchestrator__multi_comprehensive_build(
    tmpdir,
    create_orchestrator,
    docker_compose_file_multi_comprehensive_build,
    dockerfile_dummy_tests_success,
    dockerfile_dummy_tests_fail,
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
