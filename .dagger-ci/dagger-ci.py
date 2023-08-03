#!/usr/bin/python
'''
Python script to build and test Docker containers for coreboot and EDK2 compilation
'''

# Logging
# https://docs.python.org/3/howto/logging.html
# DEBUG, INFO, WARNING, ERROR, CRITICAL
import logging

import os
import sys
import subprocess
import humanize
import yaml
import argparse
import textwrap
from pprint import pformat
import anyio
import dagger

from lib.filesystem import mkdir
from lib.orchestrator import Orchestrator
from lib.results import Results


def cli(args: list = None):
    '''
    Command Line Interface
      The optional "args" is for unit-testing
    '''
    # Define CLI
    parser = argparse.ArgumentParser(
        description='Python script to build and test docker containers with dagger for firmware-action',
        epilog='https://github.com/9elements/firmware-action',
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        '-v', '--verbose',
        help='increase verbosity',
        action='store_true',
    )
    parser.add_argument(
        '-c', '--concurent',
        help='execute builds concurently',
        action='store_true',
    )
    parser.add_argument(
        '-d', '--dockerfile',
        help=textwrap.dedent('''\
                select which dockerfile to build
                - enter name from docker compose
                - multiple entries are possible
                - by default tries to build all'''),
        nargs='+',
    )
    parser.add_argument(
        '-p', '--publish',
        help='publish if build and tests succeed',
        action='store_true',
    )

    # If unit-testing parse passed 'args'
    if args is not None:
        return parser.parse_args(args=args), parser

    # If not unit-testing, parse arguments from console
    return parser.parse_args(), parser


async def main(args: list = None):
    '''
    The main function, duh
    '''
    args, parser = cli(args=args)

    # Setup nice logging
    logging.basicConfig(format='%(levelname)s: %(message)s',
                        level=logging.DEBUG if args.verbose else logging.INFO)

    # Figure out location of docker compose
    current_dir = os.path.dirname(os.path.realpath(__file__))
    repo_root_dir = current_dir
    while os.path.basename(repo_root_dir) != 'firmware-action':
        repo_root_dir = os.path.dirname(repo_root_dir)
    docker_compose_path = os.path.join(repo_root_dir, 'docker', 'compose.yaml')

    # Init the Orchestrator
    my_orchestrator = Orchestrator(
        docker_compose_path=docker_compose_path,
        concurent=args.concurent,
        publish=args.publish
    )

    # Perform builds, tests and publishing
    results = await my_orchestrator.build_test_publish(dockerfiles_override=args.dockerfile)

    # Pretty print results
    results.print()
    return results.return_code


if __name__ == '__main__':
    #####
    # Mandatory guard
    # Detailed explanation: https://stackoverflow.com/a/419185
    #####

    # LET'S DO IT !!!
    sys.exit(anyio.run(main))
