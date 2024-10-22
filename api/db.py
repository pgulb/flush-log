from passlib.hash import bcrypt
from pymongo import MongoClient
from pymongo.errors import ConnectionFailure, DuplicateKeyError


class MockClient:
    def __init__(self, not_avalaible: str):
        self.admin = self
        if not_avalaible == "true":
            self.command = self.__bad_command
        self.flush = self
        self.users = self
        if not_avalaible == "true":
            self.insert_one = self.__bad_command
        self.users_collection = {}

    def __bad_command(self, command: str):
        raise ConnectionFailure(f"Mock server not available, used command: {command}")

    def command(self, command: str) -> str:
        if command == "ping":
            return "OK"
        raise NotImplementedError

    def insert_one(self, document: dict) -> None:
        if document["_id"] not in self.users_collection.keys():
            self.users_collection[document["_id"]] = document
        else:
            raise DuplicateKeyError("Document already exists")


def hash_password(password: str) -> str:
    return bcrypt.hash(password)


def verify_pass_hash(password: str, pass_hash: str) -> bool:
    return bcrypt.verify(password, pass_hash)


def create_mongo_client(url: str) -> MongoClient:
    return MongoClient(url)


def create_mock_client(not_avalaible: str) -> MockClient:
    return MockClient(not_avalaible)
