from fastapi import status
from fastapi.testclient import TestClient

from api.main import app

client = TestClient(app)


def test_read_healthz():
    response = client.get("/healthz")
    assert response.status_code == status.HTTP_200_OK
    assert response.text == '"OK"'


def test_read_readyz():
    response = client.get("/readyz")
    assert response.status_code == status.HTTP_200_OK
    assert response.text == '"OK"'
