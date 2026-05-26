import random
import string


def random_string(size: int=10):
    return ''.join(random.choice(string.ascii_letters) for _ in range(size))