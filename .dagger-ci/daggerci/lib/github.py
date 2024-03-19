"""
Functions to deal github stuff
"""

# mypy: disable-error-code="import"

import re

from lib.env import get_env_var_value


def github_get_project_url(project_name: str = "9elements/firmware-action") -> str:
    """
    Get URL of project
    """
    return str(
        get_env_var_value(["GITHUB_SERVER_URL"], fallback="https://github.com")
        + "/"
        + get_env_var_value(["GITHUB_REPOSITORY"], fallback=project_name)
    )


def github_get_pull_request_number() -> str | None:
    """
    Return a pull request number, if exists
    else return None
    """
    ref = get_env_var_value(["GITHUB_REF"], fallback=None)
    if ref is None:
        return None
    if re.match(r"^refs\/pull\/", ref):
        return re.sub(r"^refs\/pull\/(.*)\/merge", r"\1", ref)
    return None
