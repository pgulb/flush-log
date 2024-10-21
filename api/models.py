import re

from pydantic import BaseModel, Field, validator


class User(BaseModel):
    username: str = Field(max_length=60)
    password: str = Field(min_length=8, max_length=60)

    @validator("username")
    def validate_username(cls, v):  # noqa: N805
        if not re.match(r"^[A-Za-z0-9_]+$", v):
            raise ValueError(
                "Use only alphanumeric characters and underscores for username"
            )
        return v
