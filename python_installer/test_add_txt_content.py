from add_txt_content import hash_lines, add_content, check_content


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