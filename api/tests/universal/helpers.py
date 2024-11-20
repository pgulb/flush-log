from fastapi import status
from fastapi.testclient import TestClient
from httpx import BasicAuth


def create_user(client: TestClient, username: str, password: str):
    response = client.post("/user", json={"username": username, "password": password})
    assert response.status_code == status.HTTP_201_CREATED


def create_flush(client: TestClient, username: str, password: str, flush: dict):
    response = client.put(
        "/flush",
        json=flush,
        auth=BasicAuth(username=username, password=password),
    )
    assert response.status_code == status.HTTP_201_CREATED


def create_feedback(client: TestClient, username: str, password: str, note: str):
    response = client.post(
        "/feedback",
        auth=BasicAuth(username=username, password=password),
        params={"note": note},
    )
    print(response.text)
    print(response.status_code)
    assert response.status_code == status.HTTP_201_CREATED
