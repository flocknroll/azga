import hashlib
from add_txt_content import hash_lines, add_content, check_content


def _copy_file(src, dest):
    with open(src, "rb") as s:
        with open(dest, "wb") as d:
            d.write(s.read())

def _compare_files(path1, path2):
    with open(path1, "rb") as f1:
        h1 = hashlib.md5(f1.read()).digest()
    with open(path2, "rb") as f2:
        h2 = hashlib.md5(f2.read()).digest()

    return h1 == h2


def test_hash_lines():
    data = [ "A", "B", "C" ]

    nb, _ = hash_lines(data)

    assert nb == 3


def test_check_content_found():
    found = check_content("test_data/src_found.txt", "test_data/dest.txt")
    assert found


def test_check_content_not_found():
    found = check_content("test_data/src_not_found.txt", "test_data/dest.txt")
    assert not found


def test_add_content(tmp_path):
    dest = str(tmp_path / "dest.txt")
    _copy_file("test_data/dest.txt", dest)

    add_content("test_data/src_not_found.txt", dest)

    assert _compare_files(dest, "test_data/result.txt")