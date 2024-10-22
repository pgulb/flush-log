from fastapi import status
from fastapi.testclient import TestClient


def create_user(client: TestClient, username: str, password: str):
    response = client.post("/user", json={"username": username, "password": password})
    assert response.status_code == status.HTTP_201_CREATED
