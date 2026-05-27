FROM python:3.14

WORKDIR /tyre_match

# Install system dependencies required by OpenCV
RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

COPY ../../requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
