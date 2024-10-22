import logging
import os

import fastapi
import pymongo
from db import create_mock_client, create_mongo_client, hash_password, verify_pass_hash
from fastapi import Depends, HTTPException, Response, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBasicCredentials
from httpbasic import HTTPBasic
from models import User

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

mongo_setting = os.getenv("MONGO_URL")
if mongo_setting is None:
    raise ValueError("MONGO_URL not set")
if mongo_setting == "mock":
    client = create_mock_client(os.getenv("MOCK_NOT_AVAILABLE"))
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
