FROM python:3.11-slim AS base

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

WORKDIR /app

# System deps for building some Python packages (psycopg2, etc.)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential gcc libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies first (leverage Docker layer cache)
COPY backend/req.txt /app/requirements.txt
RUN pip install --upgrade pip && pip install -r /app/requirements.txt

# Copy project
COPY backend /app

# Expose Django port
EXPOSE 8000

# Simple entrypoint: run migrations then gunicorn
CMD ["bash", "-lc", "python manage.py migrate --noinput && gunicorn backend.wsgi:application --bind 0.0.0.0:8000 --workers 3"]
