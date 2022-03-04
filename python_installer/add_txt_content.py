from ast import Bytes
import sys

from typing import Iterator, Tuple

def hash_lines(lines: Iterator[str]) -> Tuple[int, Bytes]:
    pass

def read_entire_file(path) -> Iterator[str]:
    pass

def gen_rolling_hash_file(file, lines_nb):
    lines = []

    def rolling_hash_file() -> Bytes:
        pass

    return rolling_hash_file

def add_content(src_path, dest_path):
    lines_nb, src_hash = hash_lines(read_entire_file(src_path))
    found = False

    rhf_fn = gen_rolling_hash_file(dest_path, lines_nb)

    for dst_hash in rhf_fn():
        if src_hash == dst_hash:
            found = True
            break

    if not found:
        pass


if __name__ == '__main__':
    src = sys.argv[0]
    dest = sys.argv[1]

    add_content(src, dest)