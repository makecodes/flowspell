import os

import httpx


def base_url():
    return os.getenv("FLOWSPELL_TEST_URL", "http://localhost:8266")


def client():
    url = base_url()
    return httpx.Client(base_url=url)
