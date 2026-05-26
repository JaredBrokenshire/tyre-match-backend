from utils import random_string


def test_random_string():
    res = random_string()
    assert 10 == len(res)


def test_string_size():
    res = random_string(5)
    assert 5 == len(res)