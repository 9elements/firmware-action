"""
The main class to abstract all actions
"""

# pylint: disable=too-many-instance-attributes
# mypy: disable-error-code="import"

import datetime
import logging
import os
import re
import sys
from typing import Any

try:
    import humanize
except ImportError:
    HUMANIZE_INSTALLED = False
else:
    HUMANIZE_INSTALLED = True

import anyio
import dagger
from lib.docker_compose import DockerCompose, DockerComposeValidate
from lib.env import get_env_var_value
from lib.filesystem import mkdir
from lib.git import (
    git_describe,
    git_get_branch_name,
    git_get_latest_commit_sha_long,
    git_get_latest_commit_sha_short,
    git_get_root_directory,
    git_get_tag,
)
from lib.github import github_get_project_url, github_get_pull_request_number
from lib.results import Results


class ContainerMissingTestEnvVar(Exception):
    """
    Container is missing environment variable for testing
    """


class ContainerTestFailed(Exception):
    """
    Test executed inside built container failed
    """


class Orchestrator:
    """
    The main class to abstract all actions
    """

    def __init__(
        self, docker_compose_path: str, concurrent: bool = False, publish: bool = False
    ):
        """
        There is a lot to initialize
        """
        self.docker_compose_path = docker_compose_path
        self.concurrent = concurrent

        # Build location
        #   On local system, you probably want to build in /tmp, but pipeline is different story.
        #   According to docs, Linux virtual machines have rather limited memory and disk space:
        #       7 GB of RAM
        #       14 GB of SSD space
        #   At the time of writing, all Docker containers when build have over 14 GB
        #   Docs: https://docs.github.com/en/actions/using-github-hosted-runners/about-github-hosted-runners#supported-runners-and-hardware-resources
        self.build_dir = get_env_var_value(
            ["GITHUB_WORKSPACE", "RUNNER_TMP"], fallback="/tmp"
        )
        self.cache_dir = get_env_var_value(
            ["XDG_CACHE_DIR", "RUNNER_TMP"], fallback="/tmp"
        )
        # Log location
        self.logdir = os.path.join(self.build_dir, "logs")
        mkdir(path=self.logdir)
        # GIT related info
        self.tag_sha = git_get_latest_commit_sha_long()
        self.tag_sha_short = git_get_latest_commit_sha_short()
        self.tag_tag = git_get_tag()
        self.tag_timestamp = datetime.datetime.now().strftime("%Y%m%d-%H%M%S")
        self.tag_branch = git_get_branch_name()
        self.tag_pull_request_number = github_get_pull_request_number()
        self.tag_pull_request = f"pull_request_{self.tag_pull_request_number}"
        self.git_ref = get_env_var_value(["GITHUB_REF"], fallback=git_describe())
        self.version = next(
            # Version of docker container
            # pick first non-None value from a list
            (
                i
                for i in [
                    self.tag_tag,
                    self.tag_pull_request if self.tag_pull_request_number else None,
                    self.tag_branch,
                    self.tag_sha,
                ]
                if i is not None
            ),
            None,
        )
        self.organization = "9elements"
        self.project_name = "firmware-action"
        self.project_url = github_get_project_url(
            project_name=f"{self.organization}/{self.project_name}"
        )
        # Container registry vars
        self.container_registry = get_env_var_value(
            ["env.REGISTRY", "GITHUB_REGISTRY"], None
        )
        self.container_registry_username = get_env_var_value(
            ["github.actor", "GITHUB_ACTOR"], None
        )
        self.container_registry_password = get_env_var_value(
            ["secrets.GITHUB_TOKEN", "GITHUB_TOKEN"], None
        )

        # Publishing-related info
        self.labels = {
            "org.opencontainers.image.created": datetime.datetime.now()
            .astimezone()
            .isoformat(),
            "org.opencontainers.image.licenses": "MIT",
            "org.opencontainers.image.revision": self.tag_sha,
            "org.opencontainers.image.version": self.version,
            "org.opencontainers.image.ref.name": self.version,
            "org.opencontainers.image.source": self.project_url,
            "org.opencontainers.image.url": self.project_url,
            "org.opencontainers.image.documentation": self.project_url,
        }
        # Should contain:
        # - short sha
        # - branch
        # - tag if exists
        # - latest if tag exists
        # - pull request if exists
        self.tags = [
            self.tag_sha_short,
            re.sub(r"[\/-]", r"_", self.tag_branch),
            # self.tag_timestamp,
        ]
        if self.tag_tag is not None:
            self.tags += [self.tag_tag, "latest"]
        if self.tag_pull_request_number is not None:
            self.tags += [self.tag_pull_request]

        # Publishing
        self.publish = publish
        if self.container_registry is None:
            logging.warning("Publishing: Missing container registry")
            self.publish = False
        if self.container_registry_username is None:
            logging.warning("Publishing: Missing container registry username")
            self.publish = False
        if self.container_registry_password is None:
            logging.warning("Publishing: Missing container registry token")
            self.publish = False
        if not self.publish:
            logging.warning("Skipping container publishing step")
        logging.info(
            "container registry: %s; username: %s",
            self.container_registry,
            self.container_registry_username,
        )

        # Variable(s) for storing results of builds, tests and publishing
        self.results = Results()
        self.docker_compose = DockerCompose(path=self.docker_compose_path)

    async def build_test_publish(
        self, dockerfiles_override: list[str] | None = None
    ) -> Results:
        """
        Main function to call, it will:
        - build
        - test
        - publish
        all docker containers specified in docker-compose file.
        It is also responsible for concurrent dispatching and result reporting.
        Returns results which is dictionary.
        """
        top_element = "services"
        if top_element not in self.docker_compose.get_top_elements():
            raise DockerComposeValidate(
                f'Top element "{top_element}" not found in provided docker-compose'
            )

        all_dockerfiles = self.docker_compose.get_dockerfiles(top_element=top_element)
        if dockerfiles_override is not None:
            all_dockerfiles = dockerfiles_override
        logging.info("To build: %s", all_dockerfiles)

        async with dagger.Connection(dagger.Config(log_output=sys.stderr)) as client:
            for dockerfile in all_dockerfiles:
                if self.concurrent:
                    async with anyio.create_task_group() as task_group:
                        task_group.start_soon(
                            self.__build_test_publish__,
                            client,
                            top_element,
                            dockerfile,
                        )
                else:
                    await self.__build_test_publish__(
                        client, top_element=top_element, dockerfile=dockerfile
                    )

        return self.results

    async def __build_test_publish__(
        self, client: dagger.Client, top_element: str, dockerfile: str
    ) -> None:
        """
        Build, test and publish ...
        The actual calls, logic and error handling is here.
        """
        # Prepare variables
        dockerfile_dir = os.path.join(
            os.path.dirname(self.docker_compose_path),
            self.docker_compose.yaml[top_element][dockerfile]["build"]["context"],
        )
        dockerfile_path = os.path.join(
            dockerfile_dir,
            "Dockerfile",
        )
        if not os.path.isfile(dockerfile_path):
            self.results.add(
                top_element,
                dockerfile,
                "build",
                False,
                f"File '{dockerfile_path}' not found",
            )
            return
        dockerfile_args = self.docker_compose.get_dockerfile_args(
            dockerfile=dockerfile, top_element=top_element
        )
        tarball_file = os.path.join(self.build_dir, f"{dockerfile}.tar")

        # =======
        # BUILD
        logging.info("%s/%s: BUILDING", top_element, dockerfile)
        try:
            built_docker = await self.__build__(
                client=client,
                dockerfile_dir=dockerfile_dir,
                dockerfile_args=dockerfile_args,
            )
        except dagger.ExecError as exc:
            self.results.add(top_element, dockerfile, "build", False, exc.message)
            return
        except dagger.QueryError as exc:
            self.results.add(
                top_element, dockerfile, "build", False, exc.debug_query()
            )  # type: ignore [no-untyped-call]
            return
        self.results.add(top_element, dockerfile, "build")

        # add container specific labels into self.labels
        self.labels["org.opencontainers.image.description"] = (
            f"Container for building {dockerfile}"
        )
        self.labels["org.opencontainers.image.title"] = (
            f"{self.organization}/{self.project_name}/{dockerfile}"
        )

        # add labels to the container
        for name, val in self.labels.items():
            built_docker = await built_docker.with_label(name=name, value=val)

        logging.info("Docker container labels:")
        for label in await built_docker.labels():
            logging.info("label: %s = %s", await label.name(), await label.value())

        # export as tarball
        if not await built_docker.export(tarball_file):
            self.results.add(
                top_element,
                dockerfile,
                "export",
                False,
                f"Failed to export docker container {dockerfile} as tarball",
            )
            return
        self.results.add(top_element, dockerfile, "export")

        # =======
        # TEST
        logging.info("%s/%s: TESTING", top_element, dockerfile)
        try:
            await self.__test__(client=client, tarball_file=tarball_file)
        except ContainerTestFailed:
            self.results.add(top_element, dockerfile, "test", False)
            return
        self.results.add(top_element, dockerfile, "test")

        # =======
        # PUBLISH
        if self.publish:
            logging.info("%s/%s: PUBLISHING", top_element, dockerfile)
            # pylint: disable=attribute-defined-outside-init
            self.secret_token = client.set_secret(
                "GITHUB_TOKEN", self.container_registry_password
            )
            # pylint: enable=attribute-defined-outside-init
            try:
                await self.__publish__(
                    container=built_docker,
                    dockerfile=dockerfile,
                    top_element=top_element,
                )
                self.results.add(top_element, dockerfile, "publish", True)
            except dagger.QueryError as exc:
                logging.error(exc)
                self.results.add(top_element, dockerfile, "publish", False)
                return
        else:
            self.results.add(top_element, dockerfile, "publish", False, "skip")

    async def __build__(
        self, client: dagger.Client, dockerfile_dir: str, dockerfile_args: list[Any]
    ) -> dagger.Container:
        # dockerfile_args: list[dagger.api.gen.BuildArg]) -> dagger.Container:
        # For some reason I get
        #   "AttributeError: module 'dagger' has no attribute 'api'"
        """
        Does the actual building of docker container
        """
        context_dir = client.host().directory(dockerfile_dir)
        return await context_dir.docker_build(  # type: ignore [no-any-return]
            build_args=dockerfile_args
        )

    async def __test__(self, client: dagger.Client, tarball_file: str) -> None:
        """
        Test / verify that the built container is functional by executing a script inside
        """
        # Create container from tarball
        context_tarball = client.host().file(tarball_file)
        test_container = client.pipeline("test").container().import_(context_tarball)

        # Make sure that container has environment variable "VERIFICATION_TEST"
        verification_test = await test_container.env_variable("VERIFICATION_TEST")
        if verification_test is None:
            msg = 'Container is missing "VERIFICATION_TEST" environment variable'
            logging.critical(msg)
            raise ContainerMissingTestEnvVar(msg)
        logging.debug("VERIFICATION_TEST=%s", verification_test)

        # Figure out location of test-related files
        repo_root_dir = git_get_root_directory()
        logging.debug("Repository root directory: %s", repo_root_dir)
        test_dir = os.path.join(repo_root_dir, "tests")
        logging.debug("Directory with tests: %s", test_dir)

        # Execute test
        context_test_dir = client.host().directory(test_dir)
        container_name = os.path.basename(tarball_file)
        try:
            test_container = test_container.with_directory(
                "tests", context_test_dir
            ).with_exec(
                [verification_test],
                redirect_stdout=f"{container_name}_stdout.log",
                redirect_stderr=f"{container_name}_stderr.log",
            )
            test_container = await test_container.sync()
        except dagger.ExecError as ex:
            # When command in '.with_exec()' fails, exception is raised
            #   said exception contains STDERR and STDOUT
            for std_streams in [
                [f"{container_name}_stdout.log", ex.stdout],
                [f"{container_name}_stderr.log", ex.stderr],
            ]:
                with open(
                    os.path.join(self.logdir, std_streams[0]), "w", encoding="utf-8"
                ) as logfile:
                    logfile.write(std_streams[1])
            logging.error("Test on %s failed", container_name)
            raise ContainerTestFailed(ex.message)  # pylint: disable=raise-missing-from
            # This return will execute after 'finally' completes
            #   see: https://git.sr.ht/~atomicfs/dotfiles/tree/master/item/Templates/python-except-finally-example.py
        else:
            # When command in '.with_exec()' succeeds, STDERR and STDOUT are automatically
            #   redirected into text files, which must be extracted from the container
            for std_log in [
                f"{container_name}_stdout.log",
                f"{container_name}_stderr.log",
            ]:
                await test_container.file(std_log).export(
                    os.path.join(self.logdir, std_log)
                )
            # No return here, so the execution continues normally
        finally:
            # Cleanup
            size = os.path.getsize(tarball_file)
            if HUMANIZE_INSTALLED:
                logging.info(
                    "Size of '%s' tarball is %s",
                    tarball_file,
                    humanize.naturalsize(int(size)),
                )
            else:
                logging.info("Size of '%s' tarball is %d Bytes", tarball_file, size)
            os.remove(tarball_file)

    async def __publish__(
        self, container: dagger.Container, dockerfile: str, top_element: str
    ) -> None:
        """
        Publish the built container to container registry
        """
        for tag in self.tags:
            image_ref = await container.with_registry_auth(
                address=str(self.container_registry),
                username=str(self.container_registry_username),
                secret=self.secret_token,
            ).publish(
                f"{self.container_registry}/{self.organization}/{self.project_name}/{dockerfile}:{tag}"
            )
            logging.info(
                "%s/%s: Published image to: %s", top_element, dockerfile, image_ref
            )
