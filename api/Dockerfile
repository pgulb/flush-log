FROM python:3.12.6-slim-bookworm
WORKDIR /src
COPY log_conf.yaml ./
COPY requirements.txt .
RUN pip install -r requirements.txt --no-cache-dir
COPY *.py ./
USER 1000:1000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "6789", "--log-config=log_conf.yaml"]
