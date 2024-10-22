from datetime import datetime, timedelta
from random import randint

import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.main import app

client = TestClient(app)


def test_insert_new_flush():
    test_users = {
        "testcreateflush": "asdasdasd",
        "testcreateflush2": "asdasdasd",
        "testcreateflush3": "asdasdasd",
    }
    counter = 0
    notes = [
        "w0w",
        "",
        "ęśąćź_#Ý7Ð¶;ü÷VëÏqÑõfYÛ,ÛéO¡ü$âKÜòUúæeEá2iï:ÇüN¡¹eÑ;yïùà?þ<Ù÷G¹PëgùX8ó",
    ]
    phones = [True, False, True]
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        auth = httpx.BasicAuth(username=test_user, password=test_users[test_user])
        response = client.put(
            "/flush",
            json={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (datetime.now() + timedelta(minutes=15)).isoformat(
                    timespec="minutes"
                ),
                "rating": randint(1, 10),
                "note": notes[counter],
                "phone_used": phones[counter],
            },
            auth=auth,
        )
        assert response.status_code == status.HTTP_201_CREATED
        counter += 1


def test_insert_flush_bad_auth():
    response = client.put(
        "/flush",
        json={
            "time_start": datetime.now().isoformat(timespec="minutes"),
            "time_end": (datetime.now() + timedelta(minutes=15)).isoformat(
                timespec="minutes"
            ),
            "rating": randint(1, 10),
            "note": "w0w",
            "phone_used": True,
        },
        auth=httpx.BasicAuth(
            username="nonexistent123123", password="nonexistent123123"
        ),
    )
    assert response.status_code == status.HTTP_401_UNAUTHORIZED


def test_insert_flush_time_end_before_time_start():
    create_user(client, "testflushtimeendbefore", "testflushtimeendbefore")
    response = client.put(
        "/flush",
        json={
            "time_start": datetime.now().isoformat(timespec="minutes"),
            "time_end": (datetime.now() - timedelta(minutes=15)).isoformat(
                timespec="minutes"
            ),
            "rating": randint(1, 10),
            "note": "w0w",
            "phone_used": True,
        },
        auth=httpx.BasicAuth(
            username="testflushtimeendbefore", password="testflushtimeendbefore"
        ),
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_insert_flush_rating_out_of_range():
    create_user(client, "testflushratingoutofrange", "testflushratingoutofrange")
    for rating in [-1, 0, 11, 123, -123]:
        response = client.put(
            "/flush",
            json={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (datetime.now() + timedelta(minutes=15)).isoformat(
                    timespec="minutes"
                ),
                "rating": rating,
                "note": "w0w",
                "phone_used": True,
            },
            auth=httpx.BasicAuth(
                username="testflushratingoutofrange",
                password="testflushratingoutofrange",
            ),
        )
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_insert_flush_note_too_long():
    create_user(client, "testflushnotetoolong", "testflushnotetoolong")
    response = client.put(
        "/flush",
        json={
            "time_start": datetime.now().isoformat(timespec="minutes"),
            "time_end": (datetime.now() + timedelta(minutes=15)).isoformat(
                timespec="minutes"
            ),
            "rating": randint(1, 10),
            "note": "w0w" * 100,
            "phone_used": True,
        },
        auth=httpx.BasicAuth(
            username="testflushnotetoolong",
            password="testflushnotetoolong",
        ),
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
