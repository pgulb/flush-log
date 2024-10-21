import os

from fastapi import status
from fastapi.testclient import TestClient

from ..main import app

client = TestClient(app)


def test_check_for_mongo_url():
    assert "mongodb://" in os.environ["MONGO_URL"]

def test_read_readyz():
    response = client.get("/readyz")
    assert response.status_code == status.HTTP_200_OK
    assert response.text == '"OK"'
