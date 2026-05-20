FROM python:3.14

WORKDIR /main

COPY ../../requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
