from datetime import datetime, timedelta

import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.main import app

client = TestClient(app)


def test_flush_delete():
    create_user(client, "testflushdelete", "testflushdelete")
    for i in range(1, 11):
        response = client.put(
            "/flush",
            json={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (datetime.now() + timedelta(minutes=1 * i)).isoformat(
                    timespec="minutes"
                ),
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
    for i in range(1, 11):
        response = client.request(
            "DELETE",
            "/flush",
            json={
                "time_start": datetime.now().isoformat(timespec="minutes"),
                "time_end": (datetime.now() + timedelta(minutes=1 * i)).isoformat(
                    timespec="minutes"
                ),
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
            json={
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
