import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_feedback, create_user

from api.main import app, get_feedback_count

client = TestClient(app)


def test_insert_new_feedback():
    test_users = {
        "testcreatefeedback": "asdasdasd",
        "testcreatefeedback2": "asdasdasd",
        "testcreatefeedback3": "asdasdasd",
    }
    notes = [
        "w0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0ww0w",
        "w0ww0www0ww0ww0ww0ww0ww0wdsasd",
        "ęśąćź_#Ý7Ð¶;ü÷VëÏqÑõfYÛ,ÛéO¡ü$âKÜòUúæeEá2iï:ÇüN¡¹eÑ;yïùà?þ<Ù÷G¹PëgùX8ó",
    ]
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        for i, note in enumerate(notes):
            create_feedback(
                client,
                test_user,
                test_users[test_user],
                note,
            )
            assert get_feedback_count(test_user) == i + 1


def test_insert_feedback_bad_auth():
    response = client.post(
        "/feedback",
        auth=httpx.BasicAuth(
            username="usernamenonexistent", password="passwordnexistent"
        ),
        params={"note": "notenotenotenotenotenotenotenotenotenotenotenotenote"},
    )
    assert response.status_code == status.HTTP_401_UNAUTHORIZED


def test_insert_feedback_note_too_long():
    create_user(client, "testfeedbacknotetoolong", "testfeedbacknotetoolong")
    response = client.post(
        "/feedback",
        auth=httpx.BasicAuth(
            username="testfeedbacknotetoolong", password="testfeedbacknotetoolong"
        ),
        params={
            "note": "notenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenotenote"  # noqa: E501
        },
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


def test_insert_feedback_note_too_short():
    create_user(client, "testfeedbacknotetooshort", "testfeedbacknotetooshort")
    response = client.post(
        "/feedback",
        auth=httpx.BasicAuth(
            username="testfeedbacknotetooshort", password="testfeedbacknotetooshort"
        ),
        params={"note": "note"},
    )
    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
