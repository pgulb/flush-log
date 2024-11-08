import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.main import app

client = TestClient(app)


def test_change_password():
    user = {
        "username": "testuserchangepass",
        "password": "testpass123123",
    }
    create_user(client, user["username"], user["password"])
    new_pass = "newpass123123123"
    auth = httpx.BasicAuth(username=user["username"], password=user["password"])
    result = client.put(
        "/pass_change",
        json={"username": user["username"], "password": new_pass},
        auth=auth,
    )
    assert result.status_code == status.HTTP_200_OK
    result = client.get(
        "/", auth=httpx.BasicAuth(username=user["username"], password=new_pass)
    )
    assert result.status_code == status.HTTP_200_OK
    result = client.get(
        "/", auth=httpx.BasicAuth(username=user["username"], password=user["password"])
    )
    assert result.status_code == status.HTTP_401_UNAUTHORIZED


def test_change_password_to_bad_pass():
    users = [
        {
            "username": "testuserchangepassbadpass",
            "password": "testpass123123",
            "new_pass": "newnewn",
        },
        {
            "username": "testuserchangepassbadpass2",
            "password": "testpass123123",
            "new_pass": "asdfghjklmasdfghjklmasdfghjklmasdfghjklmasdfghjklmasdfghjklmb",
        },
    ]
    for user in users:
        create_user(client, user["username"], user["password"])
        auth = httpx.BasicAuth(username=user["username"], password=user["password"])
        result = client.put(
            "/pass_change",
            json={"username": user["username"], "password": user["new_pass"]},
            auth=auth,
        )
        assert result.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
