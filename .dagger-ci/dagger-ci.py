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
import prettytable
from pprint import pformat
import anyio
import dagger

from lib.filesystem import mkdir
from lib.orchestrator import Orchestrator


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

    # =====================
    # Pretty print results
    #   Print out some summary at the end
    #   and return 1 in case any of the tasks failed
    ret_val = 0
    all_dockerfiles = []
    all_stages = []

    for top_element, te_value in results.items():
        for dockerfile, df_value in te_value.items():
            # add to list of all dockerfiles
            entry = [top_element, dockerfile]
            if entry not in all_dockerfiles:
                all_dockerfiles.append(entry)
            # check stages
            for stage, st_value in df_value.items():
                # add to list of all stages
                if stage not in all_stages:
                    all_stages.append(stage)
                if te_value[dockerfile] is False:
                    logging.error(
                        'Failed %s/%s %s stage: %s',
                        top_element,
                        dockerfile,
                        stage,
                        te_value[stage+'_msg'])
                    ret_val = 1

    # Prettytable
    my_table = prettytable.PrettyTable()
    my_table.field_names = ['container']+all_stages
    for dockerfile in all_dockerfiles:
        row = [f'{dockerfile[0]}/{dockerfile[1]}']
        for stage in all_stages:
            res = results[dockerfile[0]][dockerfile[1]]
            if stage in res:
                row.append('OK' if res[stage] else 'fail')
            else:
                row.append(res['skip'])
        my_table.add_row(row)
    print(my_table)

    return ret_val


if __name__ == '__main__':
    #####
    # Mandatory guard
    # Detailed explanation: https://stackoverflow.com/a/419185
    #####

    # LET'S DO IT !!!
    sys.exit(anyio.run(main))
