from datetime import datetime, timedelta

import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.main import app, get_flush_count

client = TestClient(app)


def test_flush_delete():
    create_user(client, "testflushdelete", "testflushdelete")
    for i in range(1, 9):
        response = client.put(
            "/flush",
            json={
                "time_start": "2012-01-19 17:00:00",
                "time_end": f"2012-01-19 17:0{i}:00",
                "rating": i,
                "note": "w0w",
                "phone_used": True,
            },
            auth=httpx.BasicAuth(
                username="testflushdelete",
                password="testflushdelete",
            ),
        )
        assert response.status_code == status.HTTP_201_CREATED
    for i in range(1, 9):
        response = client.request(
            "DELETE",
            "/flush",
            params={
                "time_start": "2012-01-19 17:00:00",
                "time_end": f"2012-01-19 17:0{i}:00",
                "rating": i,
                "note": "w0w",
                "phone_used": True,
            },
            auth=httpx.BasicAuth(
                username="testflushdelete",
                password="testflushdelete",
            ),
        )
        assert response.status_code == status.HTTP_204_NO_CONTENT


def test_flush_delete_nonexistent_flush():
    create_user(client, "testflushdeletenonexistent", "testflushdeletenonexistent")
    for i in range(1, 5):
        response = client.request(
            "DELETE",
            "/flush",
            params={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (datetime.now() + timedelta(minutes=1 * i)).isoformat(
                    timespec="minutes"
                ),
                "rating": i,
                "note": "w0w",
                "phone_used": True,
            },
            auth=httpx.BasicAuth(
                username="testflushdeletenonexistent",
                password="testflushdeletenonexistent",
            ),
        )
        assert response.status_code == status.HTTP_400_BAD_REQUEST


def test_flush_delete_byid_nonexistent_flush():
    create_user(client, "testflushdeletenonexistent2", "testflushdeletenonexistent2")
    for i in range(1, 5):
        response = client.request(
            "DELETE",
            f"/flush/{i}",
            auth=httpx.BasicAuth(
                username="testflushdeletenonexistent2",
                password="testflushdeletenonexistent2",
            ),
        )
        assert response.status_code == status.HTTP_400_BAD_REQUEST


def test_flush_delete_byid():
    create_user(client, "testflushdeletebyid", "testflushdeletebyid")
    for i in range(2):
        response = client.put(
            "/flush",
            json={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (
                    datetime.now() + timedelta(minutes=15 * (i + 1))
                ).isoformat(timespec="minutes"),
                "rating": 5,
                "note": "",
                "phone_used": False,
            },
            auth=httpx.BasicAuth(
                username="testflushdeletebyid", password="testflushdeletebyid"
            ),
        )
        assert response.status_code == status.HTTP_201_CREATED
    assert get_flush_count("testflushdeletebyid") == 2  # noqa: PLR2004
    flushes = client.get(
        "/flushes",
        auth=httpx.BasicAuth(
            username="testflushdeletebyid", password="testflushdeletebyid"
        ),
    ).json()
    for key in flushes:
        response = client.request(
            "DELETE",
            f"/flush/{key['_id']}",
            auth=httpx.BasicAuth(
                username="testflushdeletebyid",
                password="testflushdeletebyid",
            ),
        )
        assert response.status_code == status.HTTP_204_NO_CONTENT
    assert get_flush_count("testflushdeletebyid") == 0
