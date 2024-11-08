from datetime import datetime, timedelta

import pytest
from pydantic import ValidationError

from api.models import Flush

now = datetime.now()
valid_cases = [
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": (now + timedelta(minutes=15)).isoformat(timespec="minutes"),
        "rating": 1,
        "note": "w0w",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 2,
        "note": "",
        "phone_used": False,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 3,
        "note": "testtesttesttesttesttesttesttesttesttestasdasdasdasdsdasdasdattesttestasdasdasdasdsdasdasdattesttest",  # noqa: E501
        "phone_used": False,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 10,
        "note": "",
        "phone_used": True,
    },
]
invalid_cases = [
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": (now - timedelta(minutes=15)).isoformat(timespec="minutes"),
        "rating": 1,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 0,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": -1,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 11,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": 123,
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 1,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": 123,
        "rating": 1,
        "note": "",
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 1,
        "note": 1,
        "phone_used": True,
    },
    {
        "time_start": now.isoformat(timespec="minutes"),
        "time_end": now.isoformat(timespec="minutes"),
        "rating": 1,
        "note": "",
        "phone_used": "True",
    },
]


def test_valid_flush_models():
    for flush in valid_cases:
        print(flush)
        _ = Flush.model_validate(flush, strict=True)


def test_invalid_flush_models():
    for flush in invalid_cases:
        print(flush)
        with pytest.raises(ValidationError):
            _ = Flush.model_validate(flush, strict=True)
