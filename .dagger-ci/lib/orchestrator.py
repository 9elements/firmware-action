'''
The main class to abstract all actions
'''

import os
import sys
import logging
import datetime
import anyio
import dagger

from lib.filesystem import mkdir
from lib.docker_compose import DockerCompose, DockerComposeValidate
from lib.env import get_env_var_value
from lib.git import git_get_latest_commit_sha_long, git_get_latest_commit_sha_short, git_describe


class ContainerMissingTestEnvVar(Exception):
    '''
    Container is missing environment variable for testing
    '''


class ContainerTestFailed(Exception):
    '''
    Test executed inside built container failed
    '''


class Orchestrator():
    '''
    The main class to abstract all actions
    '''

    def __init__(self, docker_compose_path: str, concurent: bool = False, publish: bool = False):
        '''
        There is a lot to initialize
        '''
        self.docker_compose_path = docker_compose_path
        self.concurent = concurent

        # Build location
        #   On local system, you probably want to build in /tmp, but pipeline is different story.
        #   Accorgding to docs, Linux virtual machines have rather limited memory and disk space:
        #       7 GB of RAM
        #       14 GB of SSD space
        #   At the time of writing, all Docker containers when build have over 14 GB
        #   Docs: https://docs.github.com/en/actions/using-github-hosted-runners/about-github-hosted-runners#supported-runners-and-hardware-resources
        self.build_dir = get_env_var_value(['GITHUB_WORKSPACE', 'RUNNER_TMP'], fallback='/tmp')
        self.cache_dir = get_env_var_value(['XDG_CACHE_DIR', 'RUNNER_TMP'], fallback='/tmp')
        # Log location
        self.logdir = os.path.join(self.build_dir, 'logs')
        mkdir(path=self.logdir)
        # GIT related info
        self.tag_sha = git_get_latest_commit_sha_long()
        self.tag_sha_short = git_get_latest_commit_sha_short()
        self.git_ref = get_env_var_value(['GITHUB_REF'], fallback=git_describe())
        self.project_name = 'firmware-action'
        self.project_url = get_env_var_value(['GITHUB_SERVER_URL'], fallback='https://github.com') + \
            '/' + get_env_var_value(['GITHUB_REPOSITORY'],
                                    fallback=f'9elements/{self.project_name}')
        # Container registry vars
        self.container_registry = get_env_var_value(['env.REGISTRY'], None)
        self.container_registry_username = get_env_var_value(['github.actor'], None)
        self.container_registry_password = get_env_var_value(['secrets.GITHUB_TOKEN'], None)

        # Publishing-related info
        self.labels = {
            "org.opencontainers.image.created": datetime.datetime.now().astimezone().isoformat(),
            "org.opencontainers.image.licenses": "MIT",
            "org.opencontainers.image.revision": self.tag_sha,
            "org.opencontainers.image.source": self.project_url,
            "org.opencontainers.image.url": self.project_url,
            "org.opencontainers.image.version": "main",  # TODO
            # container specific
            "org.opencontainers.image.ref.name": "ubuntu",  # TODO
            "org.opencontainers.image.description": "Container for building coreboot_4.19",  # TODO
            # TODO
            "org.opencontainers.image.title": f"9elements/{self.project_name}/coreboot_4.19",
        }
        self.tags = [
            self.tag_sha_short,
            # TODO
            # type = schedule, pattern = {{date 'YYYYMMDD-hhmmss' tz = 'Europe/Berlin'}}
            # type = ref, event = branch
            # type = ref, event = tag
            # type = ref, event = pr
            # type = sha
        ]
        self.publish = publish
        if self.container_registry is None or self.container_registry_username is None or self.container_registry_password is None:
            logging.warning(
                'Missing environment variables present in GitHub CI, skipping container publishing step')
            self.publish = False

        # Variable(s) for storing results of builds, tests and publishing
        self.results = {}

        self.docker_compose = DockerCompose(path=self.docker_compose_path)

        # .with_label("org.opencontainers.image.title", "my-alpine")
        # .with_label("org.opencontainers.image.version", "1.0")
        # .with_label(
        #    "org.opencontainers.image.created",
        #    datetime.now(timezone.utc).isoformat(),
        # )
        # .with_label(
        #    "org.opencontainers.image.source",
        #    "https://github.com/alpinelinux/docker-alpine",
        # )
        # .with_label("org.opencontainers.image.licenses", "MIT")

    async def build_test_publish(self, dockerfiles_override: list = None):
        '''
        Main function to call, it will:
        - build
        - test
        - publish
        all docker containers specified in docker compose file.
        It is also responsible for concurrent dispatching and result reporting.
        Returns results which is dictionary.
        '''
        top_element = 'services'
        if top_element not in self.docker_compose.get_top_elements():
            raise DockerComposeValidate(
                'Top element "%s" not found in provided docker compose', top_element)

        all_dockerfiles = self.docker_compose.get_dockerfiles(top_element=top_element)
        if dockerfiles_override is not None:
            all_dockerfiles = dockerfiles_override

        async with dagger.Connection(dagger.Config(log_output=sys.stderr)) as client:
            for dockerfile in all_dockerfiles:
                if self.concurent:
                    async with anyio.create_task_group() as task_group:
                        task_group.start_soon(self.__build_test_publish__,
                                              client, top_element=top_element, dockerfile=dockerfile)
                else:
                    await self.__build_test_publish__(client, top_element=top_element, dockerfile=dockerfile)

        return self.results

    async def __build_test_publish__(self, client, top_element: str, dockerfile: str):
        '''
        Build, test and publish ...
        The actual calls, logic and error handling is here.
        '''
        # Prepare reporting
        if top_element not in self.results:
            self.results[top_element] = {}
        if dockerfile not in self.results[top_element]:
            self.results[top_element][dockerfile] = {}

        # Prepare variables
        dockerfile_dir = os.path.join(
            os.path.dirname(self.docker_compose_path),
            self.docker_compose.yaml[top_element][dockerfile]['build']['context'],
        )
        dockerfile_path = os.path.join(
            dockerfile_dir,
            'Dockerfile',
        )
        if not os.path.isfile(dockerfile_path):
            raise FileNotFoundError(dockerfile_path)
        dockerfile_args = self.docker_compose.get_dockerfile_args(
            dockerfile=dockerfile, top_element=top_element)
        tarball_file = os.path.join(self.build_dir, f"{dockerfile}.tar")

        # =======
        # BUILD
        logging.info('%s/%s: BUILDING', top_element, dockerfile)
        built_docker = await self.__build__(client=client, dockerfile_dir=dockerfile_dir, dockerfile_args=dockerfile_args)
        self.results[top_element][dockerfile]['build'] = True

        # export as tarball
        if not await built_docker.export(tarball_file):
            self.results[top_element][dockerfile]['export'] = False
            self.results[top_element][dockerfile][
                'export_msg'] = f"Failed to export docker container {dockerfile} as tarball"
            return
        self.results[top_element][dockerfile]['export'] = True

        # =======
        # TEST
        logging.info('%s/%s: TESTING', top_element, dockerfile)
        try:
            await self.__test__(client=client, tarball_file=tarball_file)
        except ContainerTestFailed as exc:
            self.results[top_element][dockerfile]['test'] = False
            return
        self.results[top_element][dockerfile]['test'] = True

        # =======
        # PUBLISH
        if self.publish:
            self.__publish__()
        else:
            self.results[top_element][dockerfile]['publish'] = False
            self.results[top_element][dockerfile]['publish_msg'] = 'skip'

        # TODO: call self.__publish__

    async def __build__(self, client, dockerfile_dir: str, dockerfile_args: list):
        '''
        Does the actual building of docker container
        '''
        context_dir = client.host().directory(dockerfile_dir)
        return await context_dir.docker_build(build_args=dockerfile_args)

    async def __test__(self, client, tarball_file: str):
        '''
        Test / verify that the built container is functional by executing a script inside
        '''
        # Create container from tarball
        context_tarball = client.host().file(tarball_file)
        test_container = (
            client.pipeline('test')
            .container()
            .import_(context_tarball)
        )

        # Make sure that container has environment variable "VERIFICATION_TEST"
        verification_test = await test_container.env_variable('VERIFICATION_TEST')
        if verification_test is None:
            msg = 'Container is missing "VERIFICATION_TEST" environment variable'
            logging.critical(msg)
            raise ContainerMissingTestEnvVar(msg)
        logging.debug('VERIFICATION_TEST=%s', verification_test)

        # Figure out location of test-related files
        current_dir = os.path.dirname(os.path.realpath(__file__))
        logging.debug('Current working directory: %s', current_dir)
        repo_root_dir = current_dir
        while os.path.basename(repo_root_dir) != self.project_name:
            repo_root_dir = os.path.dirname(repo_root_dir)
        logging.debug('Repository root directory: %s', repo_root_dir)
        test_dir = os.path.join(repo_root_dir, 'tests')
        logging.debug('Directory with tests: %s', test_dir)

        # Execute test
        context_test_dir = client.host().directory(test_dir)
        container_name = os.path.basename(tarball_file)
        try:
            test_container = (
                test_container
                .with_directory('tests', context_test_dir)
                .with_exec(
                    [verification_test],
                    redirect_stdout=f'{container_name}_stdout.log',
                    redirect_stderr=f'{container_name}_stderr.log',
                )
            )
            test_container = await test_container.sync()
        except dagger.exceptions.ExecError as ex:
            # When command in '.with_exec()' fails, exception is raised
            #   said exception contains STDERR and STDOUT
            for std in [
                [f'{container_name}_stdout.log', ex.stdout],
                [f'{container_name}_stderr.log', ex.stderr],
            ]:
                with open(os.path.join(self.logdir, std[0]), 'w', encoding='utf-8') as logfile:
                    logfile.write(std[1])
            logging.error("Test on %s failed", container_name)
            raise ContainerTestFailed(ex.message)
            # This return will execute after 'finally' completes
            #   see: https://git.sr.ht/~atomicfs/dotfiles/tree/master/item/Templates/python-except-finally-example.py
        else:
            # When command in '.with_exec()' suceeds, STDERR and STDOUT are automatically
            #   redirected into text files, which must be extracted from the container
            for std in [f'{container_name}_stdout.log', f'{container_name}_stderr.log']:
                await test_container.file(std).export(os.path.join(self.logdir, std))
            # No return here, so the execution continues normally
        finally:
            # Cleanup
            os.remove(tarball_file)

    async def __publish__(self, client, dockerfile: str, top_element: str):
        '''
        Publish the built container to container registry
        '''
        # TODO: unfinished
        try:
            registry = os.environ[ENV_VAR_CONTAINER_REGISTRY]
            username = os.environ[ENV_VAR_CONTAINER_REGISTRY_USERNAME]
            password = os.environ[ENV_VAR_CONTAINER_REGISTRY_PASSWORD]
        except KeyError:
            logging.warning(
                'Missing environment variables present in GitHub CI, skipping container publishing step')
            results[docker]['success'] = True
            results[docker]['msg'] = 'skipping publishing step'
            return

        logging.info("Publishing: %s", docker)

        await built_docker.with_registry_auth(registry, username, password).publish(f'{username}/{REPOSITORY}/{docker}')

        image_ref = docker_container[docker].publish(f'{docker}')
        logging.info('Published image to: %s', image_ref)

        # publish image to registry

        # print image address
        # print(f"Image published at: {address}")
        results[docker]['success'] = True
        results[docker]['msg'] = 'success'
        return
