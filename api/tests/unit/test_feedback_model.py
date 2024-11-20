import pytest
from models import Feedback
from pydantic import ValidationError

valid_cases = [
    "dzien dobry esaseaseas easeas ease ase as ease asedsasdasd",
    r"k;\ÎºIàðãÀv,ëm,ÓÅèïº)ÿ7dèÉó×EàrXê3liæWãÎtÅÑ%²Xse²w¥=ðÿë8+Îå6_ÊÁ¶w£,!Iaú¾%¤×øzNíæ¤æ\ØåÐo7\ÿk;\ÎºIàðãÀv,ëm,ÓÅèïº)ÿ7dèÉó×EàrXê3liæWãÎtÅÑ%²Xse²w¥=ðÿë8+Îå6_ÊÁ¶w£,!Iaú¾%¤×øzNíæ¤æ\ØåÐo7\ÿk;\ÎºIàðãÀv,ëm,ÓÅèïº)ÿ7dèÉó×EàrXê3liæWãÎtÅÑ%²Xse²w¥=ðÿë8+Îå6_ÊÁ¶w£,!Iaú¾%¤×øzNíæ¤æ\ØåÐo7\ÿ",  # noqa: RUF001
    "witam witam witam witam witamm",
    "ęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćźęśąćź",
]
invalid_cases = [
    "",
    "za krotkie",
    "za dlugieeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",  # noqa: E501
    "dalej za krotkie, dalejjjjjjj",
]


def test_valid_feedback_models():
    for feedback in valid_cases:
        print(feedback)
        _ = Feedback.model_validate({"note": feedback}, strict=True)


def test_invalid_feedback_models():
    for feedback in invalid_cases:
        print(feedback)
        with pytest.raises(ValidationError):
            _ = Feedback.model_validate({"note": feedback}, strict=True)
