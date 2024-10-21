from fastapi import status
from fastapi.testclient import TestClient

from api.main import app

client = TestClient(app)


def test_create_user():
    response = client.post("/user", json={"username": "test", "password": "testtest"})
    assert response.status_code == status.HTTP_201_CREATED


def test_create_user_fail_on_existing():
    response = client.post(
        "/user", json={"username": "existing", "password": "testtest"}
    )
    response = client.post(
        "/user", json={"username": "existing", "password": "testtest"}
    )
    assert response.status_code == status.HTTP_409_CONFLICT


def test_create_user_fail_short_password():
    response = client.post(
        "/user", json={"username": "shortpass", "password": "1234567"}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_too_long_password():
    response = client.post(
        "/user", json={"username": "shortpass", "password": "X" * 61}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_too_long_username():
    response = client.post(
        "/user", json={"username": "X" * 61, "password": "123456789"}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_bad_chars_username():
    for username in [
        "ęśąćź",
        "<script>alert('hackerman')</script>",
        ",.,1zxcasd",
    ]:
        response = client.post(
            "/user", json={"username": username, "password": "123456789"}
        )
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
