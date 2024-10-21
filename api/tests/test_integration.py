import os

from fastapi.testclient import TestClient

from ..main import app

client = TestClient(app)


def test_check_for_mongo_url():
    assert "mongodb://" in os.environ["MONGO_URL"]
