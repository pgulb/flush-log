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
    rev = flushes
    rev.reverse()
    js = response.json()["flushes"]
    for i, f in enumerate(rev):
        f["_id"] = js[i]["_id"]
    print(rev)
    print(js)
    assert js == rev


def test_getting_flushes_noflushes():
    username, password = "testgettingflushes2", "testgettingflushes2"
    create_user(client, username, password)
    response = client.get(
        "/flushes", auth=BasicAuth(username=username, password=password)
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.json()["flushes"] == []


def test_getting_flushes_with_skip():
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
    username, password = "testgettingflushesskip", "testgettingflushesskip"
    create_user(client, username, password)
    for f in flushes:
        create_flush(client, username, password, f)
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
        params={"skip": 0},
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()["flushes"]
    assert len(js) == 3
    assert js[0]["time_start"] == flushes[-1]["time_start"]
    if not response.json()["more_data_available"]:
        raise AssertionError
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
        params={"skip": 3},
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()["flushes"]
    assert len(js) == 3
    if response.json()["more_data_available"]:
        raise AssertionError
    assert js[0]["time_start"] == flushes[2]["time_start"]

    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
        params={"skip": 6},
    )
    assert response.status_code == status.HTTP_200_OK
    js = response.json()["flushes"]
    assert len(js) == 0
    if response.json()["more_data_available"]:
        raise AssertionError


def test_getting_flushes_by_id():
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
    username, password = "testgettingflushesbyid", "testgettingflushesbyid"
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


def test_getting_flushes_by_id_wrong_id():
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
    username, password = "testgettingflushesbyid2", "testgettingflushesbyid2"
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
    response = client.get(
        "/flush/randomwrongid",
        auth=BasicAuth(username=username, password=password),
    )
    assert response.status_code == status.HTTP_404_NOT_FOUND
