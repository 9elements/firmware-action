"""
Functions to help parse docker-compose
Docs: https://docs.docker.com/compose/compose-file/

!!! This class / parser is incomplete implementation of docker-compose specification !!!
It implements only specific functions that are needed in this specific project.
"""

# mypy: disable-error-code="import"

import logging
import subprocess
from pprint import pformat

import dagger
import yaml


class DockerComposeValidate(Exception):
    """
    Failed validation of docker-compose yaml file
    """


class DockerComposeMissingElement(Exception):
    """
    Failed to get element from docker-compose yaml file
    """


def select(heap: list[str], needle: str | None = None) -> str:
    """
    Select either top_element or dockerfile from compose YAML
        - as default return first defined top_element or dockerfile
        - return needle
    """
    # Default
    if needle is None:
        return heap[0]
    # If needle exists in heap, return it
    if needle in heap:
        return needle
    # If needle does not exists, raise exception
    raise ValueError


class DockerCompose:
    """
    Class to parse docker-compose file
    """

    def __init__(self, path: str):
        self.path = path
        with open(path, "r", encoding="utf-8") as composefile:
            self.yaml = yaml.safe_load(composefile.read())
            # yaml.safe_load might raise yaml.YAMLError
        self.validate()

    def validate(self) -> None:
        """
        Validate the compose.yaml file
        """
        try:
            cmd = ["docker-compose", "-f", self.path, "config"]
            output = subprocess.run(cmd, check=False, capture_output=True)
        except FileNotFoundError as exc:
            logging.error('Missing dependency "docker-compose", please install it')
            raise exc
        if output.returncode != 0:
            logging.critical('Docker-compose file "%s" failed validation', self.path)
            logging.critical(pformat(output))
            raise DockerComposeValidate("Failed docker-compose validation")

    def get_top_elements(self) -> list[str]:
        """
        Return a list of all top_elements found (list of strings)
        """
        return list(self.yaml.keys())

    def __select_top_element__(self, top_element: str | None = None) -> str:
        """
        Check if top_element in YAML
            if true, return said top_element name
            if None provided, default to the first top_element in YAML
            if false, raise ValueError exception
        """
        try:
            return select(heap=self.get_top_elements(), needle=top_element)
        except ValueError:
            raise DockerComposeMissingElement(  # pylint: disable=raise-missing-from
                f"Top element {top_element} not found in YAML file"
            )

    def get_dockerfiles(self, top_element: str | None = None) -> list[str]:
        """
        Return a list of all dockerfiles in top_element (list of strings)
        if no top_element provided, use the first one
        """
        this_top_element = self.__select_top_element__(top_element)
        return list(self.yaml[this_top_element].keys())

    def __select_dockerfile__(
        self, dockerfile: str | None = None, top_element: str | None = None
    ) -> str:
        """
        Check if dockerfile under top_element in YAML
            if true, return said dockerfile name
            if None provided, default to the first dockerfile under top_element
            if false, raise ValueError exception
        """
        this_top_element = self.__select_top_element__(top_element)
        try:
            return select(
                heap=self.get_dockerfiles(top_element=this_top_element),
                needle=dockerfile,
            )
        except ValueError:
            raise DockerComposeMissingElement(  # pylint: disable=raise-missing-from
                f"Dockerfile {dockerfile} not found in YAML file under {top_element} top_element"
            )

    def get_dockerfile_context(
        self, dockerfile: str | None = None, top_element: str | None = None
    ) -> str | None:
        """
        Return a context of given dockerfile
        """
        this_top_element = self.__select_top_element__(top_element)
        this_dockerfile = self.__select_dockerfile__(dockerfile)

        if "build" in self.yaml[this_top_element][this_dockerfile]:
            if "context" in self.yaml[this_top_element][this_dockerfile]["build"]:
                return str(
                    self.yaml[this_top_element][this_dockerfile]["build"]["context"]
                )
        return None

    def get_dockerfile_args(
        self, dockerfile: str | None = None, top_element: str | None = None
    ) -> list[dagger.BuildArg]:
        """
        Return a list of args for given dockerfile
            return list of dagger.BuildArg
            https://dagger-io.readthedocs.io/en/sdk-python-v0.8.2/client.html#dagger.BuildArg
        """
        this_top_element = self.__select_top_element__(top_element)
        this_dockerfile = self.__select_dockerfile__(dockerfile)

        if "args" in self.yaml[this_top_element][this_dockerfile]["build"]:
            return [
                dagger.BuildArg(i.split("=")[0], i.split("=")[1])
                for i in self.yaml[this_top_element][this_dockerfile]["build"]["args"]
            ]
        return []
