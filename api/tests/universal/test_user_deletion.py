import httpx
from fastapi import status
from fastapi.testclient import TestClient
from universal.helpers import create_user

from api.main import app

client = TestClient(app)


def test_delete_user_after_create():
    test_users = {
        "testdelete": "è@Ôþß6Þ3$ÆHkµw7I#BÖÛôMã5úwYí@ï2±è¬W(åa%¸:Æ¶¢Õ(WqùñÜ§ié³",  # noqa: RUF001
        "testdelete2": r"·ÀÊRÛ0ò?<Iç(-¬ÝÏ8¤m¸ZU<(÷¡ïÀô8ÄZkB>TÿßZ?|N³õÊ:\Dx~+bïY´",  # noqa: RUF001
        "testdelete3": "_|â³¿]d®²QæºãÖC¡4[$ãþÞj?O#é@×),÷µ¯ØwU¾ÃW+&Æ?Í¼7MhÔbÝAAv",  # noqa: RUF001
    }
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        auth = httpx.BasicAuth(username=test_user, password=test_users[test_user])
        response = client.delete("/user", auth=auth)
        assert response.status_code == status.HTTP_204_NO_CONTENT


def test_try_delete_with_wrong_creds():
    test_users = {
        "testdeletewrongcreds": "àûf<Èwñ.sÊ{º¸SÆ`¶ÍÔ¶±ìúª¬ÊûØ?Ìåð}ÀpÍêp©þ<ÇÝ{¬<V¥ÚöIôu3ÿï¶©Z£",  # noqa: E501, RUF001
        "testdeletewrongcreds2": r"Áð/%§À~æMè'Cô.<'p¿JNaöm}ßÖ4ÚÞRcÇ¬HÏ?À¸*f3Mèï(ÖÓwUÈø]ÛVKÕEÛxí",  # noqa: E501, RUF001
    }
    for test_user in test_users.keys():
        create_user(client, test_user, test_users[test_user])
        auth = httpx.BasicAuth(username=test_user, password="notarealpasswordlololol")
        response = client.delete("/user", auth=auth)
        assert response.status_code == status.HTTP_401_UNAUTHORIZED


def test_delete_nonexisting_user():
    test_users = {
        "testdeletewrongusername": "÷·ÐuãmÞ6Ð©L¬ý7¡ÕW¿é®c[k¾)£G¥Hæßä%¤¡Gò^T¥Vc¤ÄOçóÍtw8Ú_^=móçª",  # noqa: E501
        "testdeletewrongusername2": r"íõ;OVÓYÔN®³ï&Ö=*üqlïÎxûNãÜpüSÖ\ó£ã¼ã$Î·Åâ¦¡oëÞ¾âgýÒ´ôíl-ß/KÑ",  # noqa: E501, RUF001
    }
    for test_user in test_users.keys():
        response = client.delete(
            "/user",
            auth=httpx.BasicAuth(username=test_user, password=test_users[test_user]),
        )
        assert response.status_code == status.HTTP_401_UNAUTHORIZED
