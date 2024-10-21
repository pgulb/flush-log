import logging
import os
from random import randint

import fastapi
import pymongo
from db import create_mock_client, create_mongo_client
from fastapi import Depends, HTTPException, Response, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBasic, HTTPBasicCredentials
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


def check_creds(credentials: HTTPBasicCredentials):
    if not (credentials.username == "admin" and credentials.password == "admin"):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Basic"},
        )


@app.get("/")
def root(credentials: HTTPBasicCredentials = Depends(security)):
    check_creds(credentials)
    return f"Random string from api {randint(0, 10000)}"


@app.post("/user", status_code=status.HTTP_201_CREATED)
def create_user(user: User):
    database = client.flush
    users = database.users
    try:
        users.insert_one({"_id": user.username, "pass_hash": user.password})
    except pymongo.errors.DuplicateKeyError as e:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT, detail="User already exists"
        ) from e
    return user


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
