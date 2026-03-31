#!/usr/bin/env python3

from __future__ import annotations

import sys
from pathlib import Path


def normalize_spec(text: str) -> str:
    lines = text.splitlines()
    normalized: list[str] = []
    i = 0

    while i < len(lines):
        line = lines[i]
        stripped = line.strip()

        if stripped.startswith("openapi: 3.1"):
            indent = line[: len(line) - len(line.lstrip())]
            normalized.append(f"{indent}openapi: 3.0.3")
            i += 1
            continue

        if stripped == "type:" and i + 2 < len(lines):
            indent = line[: len(line) - len(line.lstrip())]
            item_prefix = indent + "  - "
            first = lines[i + 1]
            second = lines[i + 2]

            if first.startswith(item_prefix) and second.startswith(item_prefix):
                first_value = first[len(item_prefix) :].strip()
                second_value = second[len(item_prefix) :].strip()
                union_values = [first_value, second_value]

                if '"null"' in union_values:
                    non_null_values = [
                        value for value in union_values if value != '"null"'
                    ]
                    if len(non_null_values) == 1:
                        normalized.append(f"{indent}type: {non_null_values[0]}")
                        normalized.append(f"{indent}nullable: true")
                        i += 3
                        continue

        normalized.append(line)
        i += 1

    return "\n".join(normalized) + "\n"


def main() -> int:
    if len(sys.argv) != 3:
        print(
            "usage: prepare_openapi_for_codegen.py <input-path> <output-path>",
            file=sys.stderr,
        )
        return 1

    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])

    normalized = normalize_spec(input_path.read_text())
    output_path.write_text(normalized)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
