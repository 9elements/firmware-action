"""
Functions to deal with environment variables
"""

import os


def get_env_var_value(
    list_of_maybe_existing_env_vars: list[str], fallback: str | None
) -> str | None:
    """
    Get a value of first defined environment variable found in "list_of_maybe_existing_env_vars",
    or get "fallback".
    I often use "fallback" as 'None'
    """
    for i in list_of_maybe_existing_env_vars:
        if os.environ.get(i):
            return os.environ[i]
    return fallback
