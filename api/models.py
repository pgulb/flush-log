import re
from datetime import datetime

from pydantic import BaseModel, Field, ValidationInfo, field_validator


class User(BaseModel):
    username: str = Field(max_length=60)
    password: str = Field(min_length=8, max_length=60)

    @field_validator("username")
    @classmethod
    def validate_username(cls, v):
        if not re.match(r"^[A-Za-z0-9_]+$", v):
            raise ValueError(
                "Use only alphanumeric characters and underscores for username"
            )
        return v


class Flush(BaseModel):
    time_start: str
    time_end: str
    rating: int = Field(min=1, max=10)
    note: str = Field(min_length=0, max_length=100)
    phone_used: bool

    @field_validator("time_end")
    @classmethod
    def validate_time_end(cls, v, info: ValidationInfo):
        if datetime.fromisoformat(v) < datetime.fromisoformat(info.data["time_start"]):
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
