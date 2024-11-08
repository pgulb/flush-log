import pytest
from pydantic import ValidationError

from api.models import User

valid_cases = [
    {
        "username": "Twojastara69_PL",
        "password": "testtest",
    },
    {
        "username": "testtesttest",
        "password": "testtesttest",
    },
    {
        "username": "DzienDobryOgolem",
        "password": "testtesttesttest",
    },
    {
        "username": "testtesttesttesttest",
        "password": r"ª¢îdT)¸pEZÑ|ëî¢^ÒÍI_Êfc·wå¯>÷;g²®|l5]^ìNHÎý&À3¼M¡÷Ä¯ê{>A{¸Ð§",  # noqa: RUF001
    },
    {
        "username": "Test_tesT",
        "password": "Test_tesTęśąćźćźćźćźćź",
    },
]
invalid_cases = [
    {
        "username": "",
        "password": "testtest",
    },
    {
        "username": "testtesttest",
        "password": "",
    },
    {
        "username": "test",
        "password": "testtes",
    },
    {
        "username": "test",
        "password": "testtesttesttesttesttesttesttesttesttestasdasdasdasdsdasdasda",
    },
    {
        "username": "testęśąćź",
        "password": "testtesttest",
    },
    {
        "username": "test-test",
        "password": "testtesttest",
    },
    {
        "username": "_-_-_",
        "password": "testtesttest",
    },
    {
        "username": "ś",
        "password": "passpasspass",
    },
    {
        "username": 1123123123123,
        "password": 12323231231233,
    },
    {
        "username": ["qweqweqweewqeweqwe"],
        "password": ["qweqweqweewqeweqwe"],
    },
]


def test_valid_user_models():
    for user in valid_cases:
        print(user)
        _ = User.model_validate(user, strict=True)


def test_invalid_user_models():
    for user in invalid_cases:
        print(user)
        with pytest.raises(ValidationError):
            _ = User.model_validate(user, strict=True)
