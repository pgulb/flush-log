import csv
import io
import json

from fastapi import status
from fastapi.testclient import TestClient
from httpx import BasicAuth
from universal.helpers import create_flush, create_user

from api.main import app

client = TestClient(app)

test_flushes = [
    {
        "time_start": "2021-01-01T00:00:00",
        "time_end": "2021-01-01T01:00:00",
        "rating": 1,
        "note": "test",
        "phone_used": True,
    },
    {
        "time_start": "2021-01-01T01:00:00",
        "time_end": "2021-01-01T02:00:00",
        "rating": 2,
        "note": "teqg523qg253qgęśąśćżźźć523qg253q5g23st",
        "phone_used": False,
    },
    {
        "time_start": "2021-01-01T02:00:00",
        "time_end": "2021-01-01T03:00:00",
        "rating": 3,
        "note": "tq3dcq32gfq23g5est",
        "phone_used": False,
    },
    {
        "time_start": "2021-01-01T03:00:00",
        "time_end": "2021-01-01T04:00:00",
        "rating": 4,
        "note": "32123eweqwe",
        "phone_used": True,
    },
    {
        "time_start": "2021-01-01T04:00:00",
        "time_end": "2021-01-01T05:00:00",
        "rating": 5,
        "note": "ewewwqweqwe",
        "phone_used": False,
    },
    {
        "time_start": "2021-01-01T05:00:00",
        "time_end": "2021-01-01T06:00:00",
        "rating": 6,
        "note": "sdasdasdsds",
        "phone_used": True,
    },
]


def test_export_flushes_json():
    username, password = "testexportjson", "testexportjson"
    create_user(client, username, password)
    for f in test_flushes:
        create_flush(client, username, password, f)
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
        params={"export_format": "json"},
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.headers["Content-Type"] == "application/json"
    outfile = io.BytesIO(response.content)
    rev = test_flushes[::-1]
    outdict = json.loads(outfile.getvalue().decode("utf-8"))
    for i, f in enumerate(rev):
        for key in ["time_start", "time_end", "note", "phone_used", "rating"]:
            assert f[key] == outdict[i][key]


def test_export_flushes_csv():
    username, password = "testexportcsv", "testexportcsv"
    create_user(client, username, password)
    for f in test_flushes:
        create_flush(client, username, password, f)
    response = client.get(
        "/flushes",
        auth=BasicAuth(username=username, password=password),
        params={"export_format": "csv"},
    )
    assert response.status_code == status.HTTP_200_OK
    assert "text/csv" in response.headers["Content-Type"]
    outfile = io.BytesIO(response.content)
    outdict = list(csv.DictReader(outfile.getvalue().decode("utf-8").splitlines()))
    for i, f in enumerate(test_flushes[::-1]):
        print(f)
        print(outdict[i])
        assert f["note"] == outdict[i]["note"]
        assert f["rating"] == int(outdict[i]["rating"])
        for key in ["time_start", "time_end"]:
            assert f[key] == outdict[i][key].replace(" ", "T")
        if f["phone_used"]:
            assert outdict[i]["phone_used"] == "True"
        else:
            assert outdict[i]["phone_used"] == "False"
