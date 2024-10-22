from fastapi import status
from fastapi.testclient import TestClient

from api.db import hash_password, verify_pass_hash
from api.main import app

client = TestClient(app)


def test_create_user():
    response = client.post("/user", json={"username": "test", "password": "testtest"})
    assert response.status_code == status.HTTP_201_CREATED


def test_create_user_fail_on_existing():
    response = client.post(
        "/user", json={"username": "existing", "password": "testtest"}
    )
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
        "testhash": r"woeirujghwgr023uh4g039hr",
        "testhash2": r"jhsdfhgsdnvnvnb_!@_#>?<><>?!!@#\\\\////",
        "testhash3": r"testtest3",
        "testhash4": r"*å¤^I!]°3Ãdké=Nçß\{»gü|2Cñ²Ñ4«Ç.gÄ{\"!Áý|$ÁflÄù¢2qBáÇzLR·à(",
        "testhash5": r"vqE´Ê.;Ì³Ìøò¼Î¼/O¤ýúPly¦S3¯JkÁQ¨e*ÀC§(îN®Ä#i¶¤¼ÖWÒ",  # noqa: RUF001
    }
    for test_user in test_users.keys():
        response = client.post(
            "/user", json={"username": test_user, "password": test_users[test_user]}
        )
        assert response.status_code == status.HTTP_201_CREATED
        assert verify_pass_hash(
            test_users[test_user], hash_password(test_users[test_user])
        )


def test_verify_pass_hash_fail():
    test_pass = "tetetetetetesciwo123"
    hash_pass = hash_password(test_pass)
    response = client.post(
        "/user", json={"username": "failhash", "password": test_pass}
    )
    assert response.status_code == status.HTTP_201_CREATED
    assert not verify_pass_hash("wrongpass", hash_pass)
    assert not verify_pass_hash("alsowrong", hash_pass)
    assert not verify_pass_hash("Tetetetetetesciwo123", hash_pass)
