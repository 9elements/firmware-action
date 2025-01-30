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
    # Typical git short sha which is 8-characters long:
    # return git_get_latest_commit_sha_long()[:7]
    # 12-character long short sha:
    return git_get_latest_commit_sha_long()[:11]


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
    Return a tag if currently tag is checked-out, else return None
    """
    describe = git_describe()
    if describe is None:
        return None
    if re.match(r"^v?\d+\.\d+(?:\.\d+)?$", describe):
        return describe
    return None


def git_get_branch_name() -> str:
    """
    Return name of current branch
    """
    return (
        subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"])
        .strip()
        .decode()
    )


def git_get_root_directory() -> str:
    """
    Assuming that current working directory is part of git repository,
    get the root directory of said git repository
    """
    return (
        subprocess.check_output(["git", "rev-parse", "--show-toplevel"])
        .strip()
        .decode()
    )
