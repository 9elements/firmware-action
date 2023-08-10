#!/usr/bin/python
"""
Python script to build and test Docker containers for coreboot and EDK2 compilation
"""
# mypy: disable-error-code="import"

# Logging
# https://docs.python.org/3/howto/logging.html
# DEBUG, INFO, WARNING, ERROR, CRITICAL
import logging
import os
import sys

import anyio
from lib.cli import cli
from lib.git import git_get_root_directory
from lib.orchestrator import Orchestrator


async def main(inargs: list[str] | None = None) -> int:
    """
    The main function, duh
    """
    args = cli(args=inargs)

    # Setup nice logging
    logging.basicConfig(
        format="%(levelname)s: %(message)s",
        level=logging.DEBUG if args.verbose else logging.INFO,
    )

    # Figure out location of docker-compose
    docker_compose_path = os.path.join(
        git_get_root_directory(), "docker", "compose.yaml"
    )

    # Init the Orchestrator
    my_orchestrator = Orchestrator(
        docker_compose_path=docker_compose_path,
        concurrent=args.concurrent,
        publish=args.publish,
    )

    # Perform builds, tests and publishing
    results = await my_orchestrator.build_test_publish(
        dockerfiles_override=args.dockerfile
    )

    # Pretty print results
    results.print()
    return results.return_code  # type: ignore [no-any-return]


if __name__ == "__main__":
    #####
    # Mandatory guard
    # Detailed explanation: https://stackoverflow.com/a/419185
    #####

    # LET'S DO IT !!!
    sys.exit(anyio.run(main))
