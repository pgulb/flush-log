from datetime import datetime, timedelta
from random import randint

import httpx
from fastapi import status
from fastapi.testclient import TestClient

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
        response = client.post(
            "/user", json={"username": test_user, "password": test_users[test_user]}
        )
        assert response.status_code == status.HTTP_201_CREATED
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
