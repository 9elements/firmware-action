"""
Functions to deal interface with local git repository
"""

import re
import subprocess


def git_get_latest_commit_sha_long() -> str:
    """
    Assuming that current working directory is part of git repository,
    get the sha of latest commit and return it.
    """
    return subprocess.check_output(["git", "rev-parse", "HEAD"]).strip().decode()


def git_get_latest_commit_sha_short() -> str:
    """
    Assuming that current working directory is part of git repository,
    get the sha of latest commit and return it.
    """
    return git_get_latest_commit_sha_long()[:7]


def git_describe() -> str | None:
    """
    Run "git describe --tags"
    """
    try:
        return subprocess.check_output(["git", "describe", "--tags"]).strip().decode()
    except subprocess.CalledProcessError:
        return None


def git_get_tag() -> str | None:
    """
    Return a tag if currently tag is checked-out, else neturn None
    """
    describe = git_describe()
    if describe is None:
        return None
    if re.match(r"^v?\d+\.\d+(?:\.\d+)?$", describe):
        return describe
    return None