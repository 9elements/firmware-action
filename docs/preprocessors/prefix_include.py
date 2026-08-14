#!/usr/bin/env python3
"""mdbook preprocessor: "prefix-include".

mdBook's built-in `{{#include ...}}` helper (part of the "links"
preprocessor) replaces the line containing the directive with the raw
content of the included file. If that line carries a prefix - most
commonly blockquote markers (`> `) used for native mdBook admonitions
(`> [!TIP] ...`), but also plain indentation - only the *first* inserted
line keeps that prefix. Every subsequent line of a multi-line include
loses it, which breaks the blockquote (or indentation) for the rest of
the included content.

This preprocessor re-implements `{{#include}}` (plain includes, line
ranges and `ANCHOR`/`ANCHOR_END` regions - the same syntax documented at
https://rust-lang.github.io/mdBook/format/mdbook.html#including-files)
but re-applies the original line's prefix to *every* line of the
included content, not just the first one.

It is meant to run *before* mdBook's built-in "links" preprocessor
(see `before = ["links"]` in book.toml) so that, by the time "links"
runs, no `{{#include}}` directives are left for it to (mis-)process.
"""

import json
import re
import sys
from pathlib import Path

# Matches a whole line consisting of an optional prefix (blockquote
# markers and/or whitespace) followed by a `{{#include ...}}` directive
# and nothing else.
INCLUDE_RE = re.compile(
    r"^(?P<prefix>[ \t>]*)\{\{#include\s+(?P<spec>[^}]+?)\s*\}\}[ \t]*$",
    re.MULTILINE,
)

ANCHOR_RE = re.compile(r"ANCHOR:\s*([\w-]+)")
ANCHOR_END_RE = re.compile(r"ANCHOR_END:\s*([\w-]+)")

# A "range" spec is: N | N: | :N | N:M  (all parts optional, at least one digit or colon present)
RANGE_RE = re.compile(r"^(?P<start>\d*):(?P<end>\d*)$")


class IncludeError(Exception):
    """Raised when an `{{#include}}` directive cannot be resolved."""


def load_lines(path: Path) -> list[str]:
    """Read `path` and return its content split into lines (no newlines)."""
    try:
        text = path.read_text()
    except OSError as e:
        raise IncludeError(f"failed to read include file '{path}': {e}") from e
    # Preserve line content without trailing newline chars; splitlines()
    # drops the final empty element for a trailing newline, which is
    # what we want.
    return text.splitlines()


def strip_anchor_markers(lines: list[str]) -> list[str]:
    """Drop any `ANCHOR`/`ANCHOR_END` marker lines from `lines`."""
    return [
        line
        for line in lines
        if not ANCHOR_RE.search(line) and not ANCHOR_END_RE.search(line)
    ]


def select_by_anchor(lines: list[str], name: str) -> list[str]:
    """Return the lines between an `ANCHOR: name` / `ANCHOR_END: name` pair."""
    start_idx = None
    end_idx = None
    for i, line in enumerate(lines):
        m = ANCHOR_RE.search(line)
        if m and m.group(1) == name and start_idx is None:
            start_idx = i
            continue
        if start_idx is not None:
            m_end = ANCHOR_END_RE.search(line)
            if m_end and m_end.group(1) == name:
                end_idx = i
                break
    if start_idx is None or end_idx is None:
        raise IncludeError(f"anchor '{name}' not found")
    selected = lines[start_idx + 1 : end_idx]
    return strip_anchor_markers(selected)


def select_by_range(lines: list[str], start: str, end: str) -> list[str]:
    """Return `lines[start:end]` using mdBook's 1-indexed, inclusive range spec."""
    start_n = int(start) if start else 1
    end_n = int(end) if end else len(lines)
    # Line numbers in the spec are 1-indexed and inclusive.
    return lines[max(start_n - 1, 0) : end_n]


def resolve_include(base_dir: Path, spec: str) -> list[str]:
    """Resolve an `{{#include}}` spec (path[:line-range-or-anchor]) to lines."""
    # File paths used in this book never contain colons, so splitting
    # on the first colon is enough to separate path from the optional
    # line-range / anchor spec.
    if ":" in spec:
        rel_path, sub_spec = spec.split(":", 1)
    else:
        rel_path, sub_spec = spec, None

    file_path = (base_dir / rel_path.strip()).resolve()
    lines = load_lines(file_path)

    if sub_spec is None or sub_spec == "":
        return strip_anchor_markers(lines)

    if sub_spec.isdigit():
        return select_by_range(lines, sub_spec, sub_spec)

    m = RANGE_RE.match(sub_spec)
    if m:
        return select_by_range(lines, m.group("start"), m.group("end"))

    return select_by_anchor(lines, sub_spec)


def expand_includes(content: str, base_dir: Path) -> str:
    """Replace every `{{#include}}` directive in `content` with its resolved,
    prefix-preserved lines."""

    def repl(m: re.Match) -> str:
        prefix = m.group("prefix")
        spec = m.group("spec")
        included_lines = resolve_include(base_dir, spec)
        return "\n".join(
            f"{prefix}{line}".rstrip() if not line else f"{prefix}{line}"
            for line in included_lines
        )

    return INCLUDE_RE.sub(repl, content)


def process_chapter(chapter: dict, src_dir: Path) -> None:
    """Expand includes in-place for a single chapter's content."""
    chapter_path = chapter.get("path")
    if chapter_path is None:
        # Draft chapters have no path/content.
        return
    base_dir = (src_dir / chapter_path).parent
    try:
        chapter["content"] = expand_includes(chapter["content"], base_dir)
    except IncludeError as e:
        print(f"prefix-include: error in '{chapter_path}': {e}", file=sys.stderr)
        sys.exit(1)


def walk_items(items: list, src_dir: Path) -> None:
    """Recursively expand includes for every chapter in the book's item tree."""
    for item in items:
        if isinstance(item, dict) and "Chapter" in item:
            chapter = item["Chapter"]
            process_chapter(chapter, src_dir)
            walk_items(chapter.get("sub_items", []), src_dir)


def main() -> None:
    """Entry point: handle the `supports` subcommand, otherwise run as a
    normal mdBook preprocessor (read Book JSON from stdin, write it back to
    stdout)."""
    if len(sys.argv) > 1:
        if sys.argv[1] == "supports":
            # Supports every renderer.
            sys.exit(0)
        print(f"prefix-include: unknown argument: {sys.argv[1]}", file=sys.stderr)
        sys.exit(1)

    context, book = json.load(sys.stdin)
    root = Path(context["root"])
    src = context["config"]["book"].get("src", "src")
    src_dir = root / src

    walk_items(book["items"], src_dir)

    json.dump(book, sys.stdout)


if __name__ == "__main__":
    main()
