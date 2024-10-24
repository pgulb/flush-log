from fastapi import status
from fastapi.testclient import TestClient
from httpx import BasicAuth
from universal.helpers import create_flush, create_user

from api.main import app

client = TestClient(app)


def test_getting_flushes():
    flushes = [
        {
            "time_start": "2021-01-01T00:00:00",
            "time_end": "2021-01-01T01:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T01:00:00",
            "time_end": "2021-01-01T02:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T02:00:00",
            "time_end": "2021-01-01T03:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T03:00:00",
            "time_end": "2021-01-01T04:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T04:00:00",
            "time_end": "2021-01-01T05:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T05:00:00",
            "time_end": "2021-01-01T06:00:00",
            "rating": 5,
            "note": "test",
            "phone_used": True,
        },
    ]
    username, password = "testgettingflushes", "testgettingflushes"
    create_user(client, username, password)
    for f in flushes:
        create_flush(client, username, password, f)

    response = client.get(
        "/flushes", auth=BasicAuth(username=username, password=password)
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.json() == flushes


def test_getting_flushes_noflushes():
    username, password = "testgettingflushes2", "testgettingflushes2"
    create_user(client, username, password)
    response = client.get(
        "/flushes", auth=BasicAuth(username=username, password=password)
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.json() == []
