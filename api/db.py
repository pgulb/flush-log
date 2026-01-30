from html.parser import HTMLParser
from io import StringIO
import hashlib

import mongomock
import pymongo


def hash_password(password: str) -> str:
    return hashlib.sha512(password.encode('utf-8')).hexdigest()


class MLStripper(HTMLParser):
    def __init__(self):
        super().__init__()
        self.reset()
        self.strict = False
        self.convert_charrefs = True
        self.text = StringIO()

    def handle_data(self, d):
        self.text.write(d)

    def get_data(self):
        return self.text.getvalue()


def sanitize(note: str) -> str:
    s = MLStripper()
    s.feed(note)
    return s.get_data()


def verify_pass_hash(password: str, pass_hash: str) -> bool:
    return hashlib.sha512(password.encode('utf-8')).hexdigest() == pass_hash


def create_mongo_client(url: str) -> pymongo.MongoClient:
    return pymongo.MongoClient(url)


def create_mock_client() -> mongomock.MongoClient:
    return mongomock.MongoClient()
