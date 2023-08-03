'''
Functions to deal with CLI
'''

import argparse
import textwrap


def cli(args: list[str] | None = None) -> tuple[argparse.Namespace, argparse.ArgumentParser]:
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
