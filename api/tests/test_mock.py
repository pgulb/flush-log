import os

from fastapi.testclient import TestClient

from ..main import app

client = TestClient(app)


def test_check_for_mock():
    assert os.environ["MONGO_URL"] == "mock"
