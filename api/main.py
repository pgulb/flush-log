import csv
import datetime
import io
import json
import logging
import os
from typing import Union

import dateutil
import fastapi
import pymongo
from bson.objectid import ObjectId
from db import (
    create_mock_client,
    create_mongo_client,
    hash_password,
    sanitize,
    verify_pass_hash,
)
from fastapi import Depends, HTTPException, Query, Response, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBasicCredentials
from httpbasic import HTTPBasic
from models import Feedback, Flush, User

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


def filter_from_flush(credentials: HTTPBasicCredentials, flush: Flush) -> dict:
    return {
        "time_start": dateutil.parser.isoparse(flush.time_start),
        "time_end": dateutil.parser.isoparse(flush.time_end),
        "user_id": credentials.username,
    }


def filter_from_creds_and_id(credentials: HTTPBasicCredentials, flush_id: str) -> dict:
    return {
        "_id": ObjectId(flush_id),
        "user_id": credentials.username,
    }


def filter_from_user(credentials: HTTPBasicCredentials) -> dict:
    return {
        "_id": credentials.username,
    }


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


@app.post("/user", status_code=status.HTTP_201_CREATED)
def create_user(user: User):
    database = client.flush
    users = database.users
    pass_hash = hash_password(user.password)
    try:
        users.insert_one(
            {
                "_id": user.username,
                "pass_hash": pass_hash,
                "registered_at": datetime.datetime.now(datetime.UTC),
            }
        )
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
    flushes = database.flushes
    try:
        flushes.delete_many({"user_id": credentials.username})
        if flushes.count_documents({"user_id": credentials.username}) > 0:
            raise Exception("Error while flush deletion")
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
            filter=filter_from_flush(credentials, flush),
            update={
                "$set": {
                    "time_start": dateutil.parser.isoparse(flush.time_start),
                    "time_end": dateutil.parser.isoparse(flush.time_end),
                    "user_id": credentials.username,
                    "rating": flush.rating,
                    "note": sanitize(flush.note),
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


@app.delete("/flush", status_code=status.HTTP_204_NO_CONTENT)
def delete_flush(
    flush: Flush = Query(), credentials: HTTPBasicCredentials = Depends(security)
):
    check_creds(credentials)
    flushes = client.flush.flushes
    try:
        result = flushes.delete_one(filter=filter_from_flush(credentials, flush))
        if result.deleted_count != 1:
            raise Exception("Flush not deleted")
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error deleting flush"
        ) from e
    return flush.time_start


@app.delete("/flush/{flush_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_flush_by_id(
    flush_id: str, credentials: HTTPBasicCredentials = Depends(security)
):
    check_creds(credentials)
    flushes = client.flush.flushes
    try:
        result = flushes.delete_one(
            filter=filter_from_creds_and_id(credentials, flush_id)
        )
        if result.deleted_count != 1:
            raise Exception("Flush not deleted")
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error deleting flush"
        ) from e
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@app.get("/flushes", status_code=status.HTTP_200_OK)
def get_flushes(
    export_format: Union[str, None] = None,
    skip: Union[int, None] = None,
    credentials: HTTPBasicCredentials = Depends(security),
):
    check_creds(credentials)
    flushes = client.flush.flushes
    try:
        if skip is not None:
            entries = [
                x
                for x in flushes.find(
                    filter={
                        "user_id": credentials.username,
                    },
                    limit=3,
                    skip=skip,
                    sort=[("time_start", pymongo.DESCENDING)],
                )
            ]
        else:
            entries = [
                x
                for x in flushes.find(
                    filter={"user_id": credentials.username},
                    sort=[("time_start", pymongo.DESCENDING)],
                )
            ]
        for entry in entries:
            entry["_id"] = str(entry["_id"])
            del entry["user_id"]
        if export_format == "json":
            for e in entries:
                e["time_start"] = e["time_start"].isoformat()
                e["time_end"] = e["time_end"].isoformat()
            js = json.dumps(entries, indent=2)
            return Response(
                content=js,
                headers={"Content-Disposition": "attachment; filename=flushes.json"},
                media_type="application/json",
            )
        if export_format == "csv":
            csv_content = io.StringIO()
            writer = csv.writer(csv_content)
            if len(entries) > 0:
                writer.writerow(entries[0])
            for e in entries:
                writer.writerow(e.values())
            return Response(
                content=csv_content.getvalue(),
                headers={"Content-Disposition": "attachment; filename=flushes.csv"},
                media_type="text/csv",
            )
        return entries
    except Exception as e:
        logging.error(e)
        logging.info(type(entries))
        logging.info(entries)
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error getting flushes"
        ) from e


def get_flush_count(username: str) -> int:
    flushes = client.flush.flushes
    return flushes.count_documents({"user_id": username})


@app.put("/pass_change", status_code=status.HTTP_200_OK)
def update_password(
    user_new_pass: User, credentials: HTTPBasicCredentials = Depends(security)
):
    check_creds(credentials)
    users = client.flush.users
    try:
        result = users.update_one(
            filter=filter_from_user(credentials),
            update={"$set": {"pass_hash": hash_password(user_new_pass.password)}},
        )
        if result.matched_count == 1:
            return Response(status_code=status.HTTP_200_OK)
        raise Exception(f"matched_count == {result.matched_count}")
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error changing password"
        ) from e


@app.get("/stats", status_code=status.HTTP_200_OK)
def get_flush_stats(credentials: HTTPBasicCredentials = Depends(security)):
    check_creds(credentials)
    flushes = client.flush.flushes
    try:
        result = flushes.aggregate(
            [
                {"$match": {"user_id": credentials.username}},
                {
                    "$group": {
                        "_id": "$user_id",
                        "flushCount": {"$sum": 1},
                        "totalTime": {
                            "$sum": {
                                "$dateDiff": {
                                    "startDate": "$time_start",
                                    "endDate": "$time_end",
                                    "unit": "minute",
                                }
                            }
                        },
                        "meanTime": {
                            "$avg": {
                                "$dateDiff": {
                                    "startDate": "$time_start",
                                    "endDate": "$time_end",
                                    "unit": "minute",
                                }
                            }
                        },
                        "meanRating": {"$avg": "$rating"},
                        "phoneUsedCount": {"$sum": {"$cond": ["$phone_used", 1, 0]}},
                    }
                },
                {
                    "$addFields": {
                        "percentPhoneUsed": {
                            "$multiply": [
                                {"$divide": ["$phoneUsedCount", "$flushCount"]},
                                100,
                            ]
                        }
                    }
                },
            ]
        )
        json_stats = result.to_list()
        if json_stats == []:
            return {
                "flushCount": 0,
                "totalTime": 0,
                "meanTime": 0,
                "meanRating": 0,
                "phoneUsedCount": 0,
                "percentPhoneUsed": 0,
            }
        json_stats = json_stats[0]
        json_stats["meanRating"] = int(json_stats["meanRating"])
        json_stats["percentPhoneUsed"] = int(json_stats["percentPhoneUsed"])
        del json_stats["_id"]
        return json_stats
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error getting stats"
        ) from e


@app.post("/feedback", status_code=status.HTTP_201_CREATED)
def give_feedback(
    feedback: Feedback = Query(), credentials: HTTPBasicCredentials = Depends(security)
):
    check_creds(credentials)
    feedbacks = client.flush.feedbacks
    try:
        feedbacks.insert_one(
            {
                "user_id": credentials.username,
                "note": feedback.note,
                "submission_time": datetime.datetime.now(datetime.UTC),
            }
        )
        return Response(status_code=status.HTTP_201_CREATED)
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Error giving feedback"
        ) from e


def get_feedback_count(username: str) -> int:
    feedbacks = client.flush.feedbacks
    return feedbacks.count_documents({"user_id": username})
