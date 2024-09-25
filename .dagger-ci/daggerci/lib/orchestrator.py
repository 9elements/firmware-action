"""
The main class to abstract all actions
"""

# pylint: disable=too-many-instance-attributes
# mypy: disable-error-code="import"

import datetime
import logging
import os
import platform
import re
import sys
from typing import Any

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


def get_current_arch() -> str:
    """
    Get CPU architecture of current machine
    Examples: 'amd64' or 'arm64'
    """
    # Docs: https://docs.python.org/3/library/platform.html#platform.machine
    current_arch = platform.machine().lower()
    arch_dict = {
        "x86_64": "amd64",
        "aarch64": "arm64",
    }
    return arch_dict[current_arch]


def get_current_platform() -> str:
    """
    Get platform string of current machine
    Examples: 'linux/amd64' or 'windows/arm64'
    """
    # Figure out Operating System, aka Platform
    # Docs: https://docs.python.org/3/library/sys.html#sys.platform
    current_platform = sys.platform.lower()
    platform_dict = {
        "win32": "windows",
        "cygwin": "windows",
    }

    # Figure out CPU
    current_arch = get_current_arch()

    # pylint: disable=consider-using-f-string
    return "{}/{}".format(platform_dict[current_platform], current_arch)


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
                    await self.__build_test_publish__(client, top_element, dockerfile)

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

        # =======
        # BUILD
        logging.info("%s/%s: BUILDING", top_element, dockerfile)
        variants = await self.__build__(
            client=client,
            dockerfile=dockerfile,
            dockerfile_dir=dockerfile_dir,
            dockerfile_args=dockerfile_args,
            top_element=top_element,
        )
        if not variants:
            return

        # add container specific labels into self.labels
        self.labels["org.opencontainers.image.description"] = (
            f"Container for building {dockerfile}"
        )
        self.labels["org.opencontainers.image.title"] = (
            f"{self.organization}/{self.project_name}/{dockerfile}"
        )

        # add labels to the container
        for key, _ in variants.items():
            for name, val in self.labels.items():
                variants[key] = await variants[key].with_label(name=name, value=val)

        logging.info("Docker container labels:")
        for key, _ in variants.items():
            for label in await variants[key].labels():
                logging.info("label: %s = %s", await label.name(), await label.value())

        # =======
        # TEST
        logging.info("%s/%s: TESTING", top_element, dockerfile)
        try:
            await self.__test__(
                client=client,
                test_container=variants[get_current_arch()],
                test_container_name=dockerfile,
            )
        except ContainerTestFailed:
            self.results.add(
                top_element, dockerfile, f"test {get_current_arch()}", False
            )
            return
        self.results.add(top_element, dockerfile, f"test {get_current_arch()}")

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
                    variants=list(variants.values()),
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
        self,
        client: dagger.Client,
        dockerfile: str,
        dockerfile_dir: str,
        dockerfile_args: list[Any],
        top_element: str,
    ) -> dict[str, dagger.Container]:
        # dockerfile_args: list[dagger.api.gen.BuildArg]) -> dagger.Container:
        # For some reason I get
        #   "AttributeError: module 'dagger' has no attribute 'api'"

        # pylint: disable=too-many-arguments
        """
        Does the actual building of docker container
        """

        # Initially I wanted to use our existing setup to build all wanted platforms, but unfortunately
        #   there are issues with native emulation when building tool-chains for coreboot and edk2.
        # For that reason we cannot build the multi-arch container as shown in cookbook as is
        #   https://docs.dagger.io/cookbook/#build-multi-arch-image
        # We have to compile the toolchains separately on the architecture for which the toolchain
        #   is intended to run on, and then copy it into the container. This requires changes to the Dockerfile
        #   and it means that building container locally will be much more complicated.

        platforms = [
            "amd64",
            "arm64",
        ]
        context_dir = client.host().directory(dockerfile_dir)
        platform_variants = {}

        for p in platforms:
            try:
                logging.info("** building platform: %s", p)
                container = await context_dir.docker_build(  # type: ignore [no-any-return]
                    platform=dagger.Platform("linux/"+p),
                    build_args=dockerfile_args + [dagger.BuildArg("TARGETARCH", p)],
                )
                platform_variants[p] = container
            except dagger.ExecError as exc:
                logging.error("Dagger execution error")
                self.results.add(
                    top_element, dockerfile, f"build {p}", False, exc.message
                )
                return {}
            except dagger.QueryError as exc:
                logging.error(
                    "Dagger query error, try this: https://archive.docs.dagger.io/0.9/235290/troubleshooting/#dagger-pipeline-is-unable-to-resolve-host-names-after-network-configuration-changes"
                )
                self.results.add(
                    top_element,
                    dockerfile,
                    f"build {p}",
                    False,
                    exc.debug_query(),
                )  # type: ignore [no-untyped-call]
                return {}
            self.results.add(top_element, dockerfile, f"build {p}")

        return platform_variants

    async def __test__(
        self,
        client: dagger.Client,
        test_container: dagger.Container,
        test_container_name: str,
    ) -> None:
        """
        Test / verify that the built container is functional by executing a script inside
        """
        # pylint: disable=too-many-locals

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
        try:
            test_container = test_container.with_directory(
                "tests", context_test_dir
            ).with_exec(
                [verification_test],
                redirect_stdout=f"{test_container_name}_stdout.log",
                redirect_stderr=f"{test_container_name}_stderr.log",
            )
            test_container = await test_container.sync()
        except dagger.ExecError as ex:
            # When command in '.with_exec()' fails, exception is raised
            #   said exception contains STDERR and STDOUT
            for std_streams in [
                [f"{test_container_name}_stdout.log", ex.stdout],
                [f"{test_container_name}_stderr.log", ex.stderr],
            ]:
                with open(
                    os.path.join(self.logdir, std_streams[0]), "w", encoding="utf-8"
                ) as logfile:
                    logfile.write(std_streams[1])
            logging.error("Test on %s failed", test_container_name)
            raise ContainerTestFailed(ex.message)  # pylint: disable=raise-missing-from
            # This return will execute after 'finally' completes
            #   see: https://git.sr.ht/~atomicfs/dotfiles/tree/master/item/Templates/python-except-finally-example.py

        # When command in '.with_exec()' succeeds, STDERR and STDOUT are automatically
        #   redirected into text files, which must be extracted from the container
        for std_log in [
            f"{test_container_name}_stdout.log",
            f"{test_container_name}_stderr.log",
        ]:
            await test_container.file(std_log).export(
                os.path.join(self.logdir, std_log)
            )
        # No return here, so the execution continues normally

    async def __publish__(
        self,
        variants: list[dagger.Container],
        dockerfile: str,
        top_element: str,
    ) -> None:
        """
        Publish the built container to container registry
        """

        # Get the first container and use it as base
        container = variants.pop(1)

        for tag in self.tags:
            image_ref = await container.with_registry_auth(
                address=str(self.container_registry),
                username=str(self.container_registry_username),
                secret=self.secret_token,
            ).publish(
                f"{self.container_registry}/{self.organization}/{self.project_name}/{dockerfile}:{tag}",
                # add remaining containers:
                platform_variants=variants,
            )
            logging.info(
                "%s/%s: Published image to: %s", top_element, dockerfile, image_ref
            )
