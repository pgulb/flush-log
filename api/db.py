from pymongo import MongoClient
from pymongo.errors import ConnectionFailure


class MockClient:
    def __init__(self, not_avalaible: str):
        self.admin = self
        if not_avalaible == "true":
            self.command = self.__bad_command

    def __bad_command(self, command: str):
        raise ConnectionFailure(f"Mock server not available, used command: {command}")

    def command(self, command: str) -> str:
        if command == "ping":
            return "OK"
        raise NotImplementedError


def create_mongo_client(url: str) -> MongoClient:
    return MongoClient(url)


def create_mock_client(not_avalaible: str) -> MockClient:
    return MockClient(not_avalaible)
