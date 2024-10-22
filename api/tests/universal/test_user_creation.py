import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.db import hash_password, verify_pass_hash
from api.main import app

client = TestClient(app)


def test_create_user():
    create_user(client, "test", "testtest")


def test_create_user_fail_on_existing():
    create_user(client, "existing", "testtest")
    response = client.post(
        "/user", json={"username": "existing", "password": "testtest"}
    )
    assert response.status_code == status.HTTP_409_CONFLICT


def test_create_user_fail_short_password():
    response = client.post(
        "/user", json={"username": "shortpass", "password": "1234567"}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_too_long_password():
    response = client.post(
        "/user", json={"username": "shortpass", "password": "X" * 61}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_too_long_username():
    response = client.post(
        "/user", json={"username": "X" * 61, "password": "123456789"}
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_create_user_fail_bad_chars_username():
    for username in [
        "ęśąćź",
        "<script>alert('hackerman')</script>",
        ",.,1zxcasd",
    ]:
        response = client.post(
            "/user", json={"username": username, "password": "123456789"}
        )
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_verify_pass_hash_after_create():
    test_users = {
        "testhash": "woeirujghwgr023uh4g039hr",
        "testhash2": "jhsdfhgsdnvnvnb_!@_#>?<><>?!!@#\\\\////",
        "testhash3": "testtest3",
        "testhash4": r'*å¤^I!]°3Ãdké=Nçß\{»gü|2Cñ²Ñ4«Ç.gÄ{"!Áý|$ÁflÄù¢2qBáÇzLR·à(',
        "testhash5": "vqE´Ê.;Ì³Ìøò¼Î¼/O¤ýúPly¦S3¯JkÁQ¨e*ÀC§(îN®Ä#i¶¤¼ÖWÒ",  # noqa: RUF001
    }
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        assert verify_pass_hash(
            test_users[test_user], hash_password(test_users[test_user])
        )


def test_verify_pass_hash_fail():
    test_pass = "tetetetetetesciwo123"
    hash_pass = hash_password(test_pass)
    create_user(client, "failhash", test_pass)
    assert not verify_pass_hash("wrongpass", hash_pass)
    assert not verify_pass_hash("alsowrong", hash_pass)
    assert not verify_pass_hash("Tetetetetetesciwo123", hash_pass)


def test_getting_root_as_user():
    test_users = {
        "testlogin": '§hÊÑZ¥B=G¢lýi¦Ã¨,ÞõWå®n/½ÚN"Z-Ô<',
        "testlogin2": "ÓÉÑéúíòæ,&qjá'géø¦½2¹¢ÎQöIgK?¢G{",
        "testlogin3": "Z3Þ<ìla%f'bYdëlwÌNy^øíó»6äujgF©·",
    }
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        auth = httpx.BasicAuth(username=test_user, password=test_users[test_user])
        response = client.get("/", auth=auth)
        assert response.status_code == status.HTTP_200_OK
        assert response.text == f'"Hello {test_user}!"'
