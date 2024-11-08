import re
from datetime import datetime

from pydantic import BaseModel, ValidationInfo, field_validator


class User(BaseModel):
    username: str
    password: str

    @field_validator("username")
    @classmethod
    def validate_username(cls, v):
        if not re.match(r"^[A-Za-z0-9_]+$", v):
            raise ValueError(
                "Use only alphanumeric characters and underscores for username"
            )
        if len(v) > 60:  # noqa: PLR2004
            raise ValueError("Username must be at most 60 characters long")
        return v

    @field_validator("password")
    @classmethod
    def validate_password(cls, v):
        if len(v) < 8 or len(v) > 60:  # noqa: PLR2004
            raise ValueError("Password must be 8-60 characters long")
        return v


class Flush(BaseModel):
    time_start: str
    time_end: str
    rating: int
    note: str
    phone_used: bool

    @field_validator("time_end")
    @classmethod
    def validate_time_end(cls, v, info: ValidationInfo):
        try:
            time_start = info.data["time_start"]
        except KeyError as e:
            raise ValueError("Probably wrong time format") from e
        if datetime.fromisoformat(v) < datetime.fromisoformat(time_start):
            raise ValueError("End time must be after start time")
        return v

    @field_validator("time_start", "time_end")
    @classmethod
    def validate_time(cls, v):
        try:
            datetime.fromisoformat(v)
        except Exception as e:
            raise ValueError("Invalid datetime format") from e
        return v

    @field_validator("rating")
    @classmethod
    def validate_rating(cls, v):
        if v < 1 or v > 10:  # noqa: PLR2004
            raise ValueError("Rating must be 1-10")
        return v

    @field_validator("note")
    @classmethod
    def validate_note(cls, v):
        if len(v) > 100:  # noqa: PLR2004
            raise ValueError("Note must be at most 100 characters")
        return v
