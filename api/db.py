import mongomock
import pymongo
from passlib.hash import bcrypt


def hash_password(password: str) -> str:
    return bcrypt.hash(password)


def verify_pass_hash(password: str, pass_hash: str) -> bool:
    return bcrypt.verify(password, pass_hash)


def create_mongo_client(url: str) -> pymongo.MongoClient:
    return pymongo.MongoClient(url)


def create_mock_client() -> mongomock.MongoClient:
    return mongomock.MongoClient()
