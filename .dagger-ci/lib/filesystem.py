'''
Filesystem-related functions
'''

import os


def mkdir(path: str) -> None:
    '''
    Equivalent to "mkdir -p"
    '''
    if not os.path.isdir(path):
        os.makedirs(path, exist_ok=True)
