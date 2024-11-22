from fastapi import status
from fastapi.testclient import TestClient
from httpx import BasicAuth
from universal.helpers import create_flush, create_user

from api.main import app

client = TestClient(app)


def test_flush_edit():
    flushes = [
        {
            "time_start": "2021-01-01T00:00:00",
            "time_end": "2021-01-01T01:00:00",
            "rating": 5,
            "note": "dddd",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T01:00:00",
            "time_end": "2021-01-01T02:00:00",
            "rating": 8,
            "note": "test",
            "phone_used": True,
        },
        {
            "time_start": "2021-01-01T02:00:00",
            "time_end": "2021-01-01T03:00:00",
            "rating": 9,
            "note": "essa",
            "phone_used": False,
        },
    ]
    username, password = "testedittingflushes", "testedittingflushes"
    create_user(client, username, password)
    for f in flushes:
        create_flush(client, username, password, f)
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()["flushes"]
    js.reverse()
    assert len(js) == len(flushes)
    response = client.put(
        f"/flush/{js[0]['_id']}",
        auth=BasicAuth(username=username, password=password),
        params={
            "time_start": "2021-01-01T02:00:00",
            "time_end": "2021-01-01T03:00:00",
            "rating": 10,  # changed
            "note": "essa+edit",  # changed
            "phone_used": False,
        },
    )
    assert response.status_code == status.HTTP_200_OK
    response = client.put(
        f"/flush/{js[2]['_id']}",
        auth=BasicAuth(username=username, password=password),
        params={
            "time_start": "2021-01-01T01:00:00",  # changed
            "time_end": "2021-01-01T01:10:00",  # changed
            "rating": 5,
            "note": "dddd",
            "phone_used": False,  # changed
        },
    )
    assert response.status_code == status.HTTP_200_OK
    flushes[2]["rating"] = 10
    flushes[2]["note"] = "essa+edit"
    flushes[0]["time_start"] = "2021-01-01T01:00:00"
    flushes[0]["time_end"] = "2021-01-01T01:10:00"
    flushes[0]["phone_used"] = False
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()["flushes"]
    js.reverse()
    assert len(js) == len(flushes)
    for i, f in enumerate(js):
        response = client.get(
            f"/flush/{f['_id']}", auth=BasicAuth(username=username, password=password)
        )
        assert response.status_code == status.HTTP_200_OK
        f_js = response.json()
        assert f_js["time_start"] == flushes[i]["time_start"]
        assert f_js["time_end"] == flushes[i]["time_end"]
        assert f_js["rating"] == flushes[i]["rating"]
        assert f_js["note"] == flushes[i]["note"]
        assert f_js["phone_used"] == flushes[i]["phone_used"]
    fake_id = "fakeid"
    response = client.put(
        f"/flush/{fake_id}",
        auth=BasicAuth(username=username, password=password),
        params={
            "time_start": "2021-01-01T01:00:00",
            "time_end": "2021-01-01T01:10:00",
            "rating": 5,
            "note": "dddd",
            "phone_used": False,
        },
    )
    assert response.status_code == status.HTTP_400_BAD_REQUEST
