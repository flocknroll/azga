import sys
import hashlib

from typing import Iterable, Tuple


def hash_lines(lines: Iterable[str]) -> Tuple[int, bytes]:
    """
    Returns a MD5 digest of the lines.
    Strips the new line character if present.
    """
    md5 = hashlib.md5()
    cnt = 0
    for line in lines:
        cnt += 1
        line = line.rstrip("\n")
        md5.update(line.encode("UTF-8"))

    return (cnt, md5.digest())


def read_entire_file(path) -> Iterable[str]:
    """
    Reads all the lines from a file
    """
    with open(path, "rt") as f:
        lines = f.readlines()

    return lines


def rolling_hash_file(path, lines_nb) -> Iterable[bytes]:
    """
    Iterate through a file and returns the MD5 digests of the lines grouped by the specified number.
    """
    lines = []
    with open(path, "rt") as f:
        for line in f.readlines():
            lines.append(line)
            if len(lines) > lines_nb:
                lines = lines[1:]

            if len(lines) == lines_nb:
                _, hash = hash_lines(lines)
                yield hash


def check_content(src_path, dest_path):
    """
    Checks if the source file content if present in the destination file.
    """
    lines_nb, src_hash = hash_lines(read_entire_file(src_path))
    found = False

    for dest_hash in rolling_hash_file(dest_path, lines_nb):
        if src_hash == dest_hash:
            found = True
            break

    return found


def add_content(src_path, dest_path):
    """
    Add the content of the source file at the end of the destination file if not already present.
    """
    found = check_content(src_path, dest_path)

    if not found:
        print("Not found - appending data")
        src = read_entire_file(src_path)
        with open(dest_path, "at") as f:
            f.write("\n")
            f.writelines(src)
    else:
        print("Data found")


if __name__ == '__main__':
    src = sys.argv[1]
    dest = sys.argv[2]

    add_content(src, dest)