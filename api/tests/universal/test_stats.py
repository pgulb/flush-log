import os

import pytest
from fastapi import status
from fastapi.testclient import TestClient
from httpx import BasicAuth
from universal.helpers import create_flush, create_user

from api.main import app

if os.environ["MONGO_URL"] == "mock":
    pytest.skip("Skipping stats tests on mock db", allow_module_level=True)
client = TestClient(app)


def test_getting_flush_stats():
    flushes = [
        {
            "time_start": "2021-01-01T00:00:00",
            "time_end": "2021-01-01T01:00:00",
            "rating": 1,
            "note": "test",
            "phone_used": False,
        },
        {
            "time_start": "2021-01-01T01:00:00",
            "time_end": "2021-01-01T02:00:00",
            "rating": 2,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T02:00:00",
            "time_end": "2021-01-01T03:00:00",
            "rating": 3,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T03:00:00",
            "time_end": "2021-01-01T04:00:00",
            "rating": 4,
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
            "rating": 6,
            "note": "test",
            "phone_used": True,
        },
    ]
    username, password = "teststats", "teststats"
    create_user(client, username, password)
    for f in flushes:
        create_flush(client, username, password, f)

    response = client.get(
        "/stats", auth=BasicAuth(username=username, password=password)
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()
    assert js["flushCount"] == 6
    assert js["totalTime"] == 360
    assert js["meanTime"] == 60
    assert js["meanRating"] == 3
    assert js["phoneUsedCount"] == 5
    assert js["percentPhoneUsed"] == 83


def test_getting_flush_stats_noflushes():
    username, password = "teststatsempty", "teststatsempty"
    create_user(client, username, password)

    response = client.get(
        "/stats", auth=BasicAuth(username=username, password=password)
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()
    assert js["flushCount"] == 0
    assert js["totalTime"] == 0
    assert js["meanTime"] == 0
    assert js["meanRating"] == 0
    assert js["phoneUsedCount"] == 0
    assert js["percentPhoneUsed"] == 0
