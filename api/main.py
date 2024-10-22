import logging
import os

import dateutil
import fastapi
import pymongo
from db import create_mock_client, create_mongo_client, hash_password, verify_pass_hash
from fastapi import Depends, HTTPException, Response, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBasicCredentials
from httpbasic import HTTPBasic
from models import Flush, User

app = fastapi.FastAPI()
origins = [
    "*",
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
security = HTTPBasic()

logging.basicConfig(level=logging.INFO, format="%(levelname)s %(asctime)s %(message)s")
logging.getLogger("passlib").setLevel(logging.ERROR)

mongo_setting = os.getenv("MONGO_URL")
if mongo_setting is None:
    raise ValueError("MONGO_URL not set")
if mongo_setting == "mock":
    client = create_mock_client()
    logging.info("Using mock client")
else:
    client = create_mongo_client(mongo_setting)
    logging.info("Using mongo client")
    if "/?" in mongo_setting:
        logging.info(f"client options: {mongo_setting.split('/?')[1]}")


def raise_basic_exception():
    raise HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Incorrect username or password",
        headers={"WWW-Authenticate": "Basic"},
    )


def check_creds(credentials: HTTPBasicCredentials):
    database = client.flush
    users = database.users
    user = users.find_one({"_id": credentials.username})
    if user is None:
        raise_basic_exception()
    if not verify_pass_hash(credentials.password.encode("utf-8"), user["pass_hash"]):
        raise_basic_exception()


@app.get("/")
def root(credentials: HTTPBasicCredentials = Depends(security)):
    check_creds(credentials)
    return f"Hello {credentials.username}!"


@app.post("/user", status_code=status.HTTP_201_CREATED)
def create_user(user: User):
    database = client.flush
    users = database.users
    pass_hash = hash_password(user.password)
    try:
        users.insert_one({"_id": user.username, "pass_hash": pass_hash})
    except pymongo.errors.DuplicateKeyError as e:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT, detail="User already exists"
        ) from e
    return user.username


@app.delete("/user", status_code=status.HTTP_204_NO_CONTENT)
def delete_user(credentials: HTTPBasicCredentials = Depends(security)):
    check_creds(credentials)
    database = client.flush
    users = database.users
    try:
        result = users.delete_one({"_id": credentials.username})
        if result.deleted_count != 1:
            raise Exception("User not deleted")
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Error while deleting account",
        ) from e
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@app.put("/flush", status_code=status.HTTP_201_CREATED)
def create_update_flush(
    flush: Flush, credentials: HTTPBasicCredentials = Depends(security)
):
    check_creds(credentials)
    flushes = client.flush.flushes
    try:
        result = flushes.update_one(
            filter={
                "time_start": dateutil.parser.isoparse(flush.time_start),
                "time_end": dateutil.parser.isoparse(flush.time_end),
                "user_id": credentials.username,
            },
            update={
                "$set": {
                    "time_start": dateutil.parser.isoparse(flush.time_start),
                    "time_end": dateutil.parser.isoparse(flush.time_end),
                    "user_id": credentials.username,
                    "rating": flush.rating,
                    "note": flush.note,
                    "phone_used": flush.phone_used,
                }
            },
            upsert=True,
        )
        if result.matched_count == 1:
            return Response(flush.time_start, status_code=status.HTTP_200_OK)
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error adding flush"
        ) from e
    return flush.time_start


@app.get("/healthz", status_code=status.HTTP_200_OK)
def healthz():
    return "OK"


@app.get("/readyz", status_code=status.HTTP_200_OK)
def readyz():
    try:
        client.admin.command("ping")
    except Exception:
        return Response("NOT OK", status_code=status.HTTP_503_SERVICE_UNAVAILABLE)
    return "OK"
