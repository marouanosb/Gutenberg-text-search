# Docker + VPS Deployment Guide

This sets up a two-container stack:
- `backend`: Django + Gunicorn on port 8000
- `frontend`: Nginx serving the React build, proxying `/server/` to the backend

Both are wired via `docker-compose.yml` and run on the same Docker network. The frontend listens on port `80` on the host by default.

## 1) Prerequisites on your VPS
- A fresh VPS (Ubuntu 22.04+ recommended)
- A user with sudo
- Domain (optional but recommended) pointing to your VPS IP

Install Docker and Compose:
```bash
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo $VERSION_CODENAME) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Verify:
```bash
sudo docker --version
sudo docker compose version
```

## 2) Get your code onto the VPS
Option A: Git clone (if repo is public or you have credentials):
```bash
git clone https://github.com/marouanosb/Gutenberg-text-search.git
cd Gutenberg-text-search
```

Option B: Copy your local folder (Windows PowerShell example):
```powershell
# From your local machine
# Replace user@server with your VPS SSH
scp -r "c:\\Users\\PC\\OneDrive\\Bureau\\Gutenberg-text-search" user@server:~
```

## 3) One-time code changes you should make
- Frontend already uses `VITE_API_BASE` with default `"/server/books/?"` (good for Nginx proxy).
- In Django `backend/backend/settings.py`, add your domain/IP to `ALLOWED_HOSTS` (e.g. `['yourdomain.com', 'your.ip.addr', 'localhost']`). For quick tests you can set `ALLOWED_HOSTS = ['*']` but prefer explicit hosts in production.

If you change settings, commit them then pull/copy to the VPS.

## 4) Build and run on the VPS
From the repo root (where `docker-compose.yml` is):
```bash
sudo docker compose build
sudo docker compose up -d
```

This starts:
- `gutenberg-backend` on internal port 8000 (published to host 8000)
- `gutenberg-frontend` on host port 80

Check containers:
```bash
sudo docker compose ps
sudo docker compose logs -f backend
sudo docker compose logs -f frontend
```

Visit `http://YOUR_SERVER_IP/` — the app should load. API is proxied at `http://YOUR_SERVER_IP/server/books/`.

## 5) Data and volumes
- SQLite DB is bind-mounted: `./backend/db.sqlite3 -> /app/db.sqlite3` so data persists on the VPS.
- `./books` is mounted read-only into the backend at `/app/books`.

## 6) TLS (optional but recommended)
Simplest: put a reverse proxy like Caddy or Nginx Proxy Manager on the VPS in front of this stack to terminate HTTPS for your domain and forward to `http://127.0.0.1:80`.

## 7) Common operations
Rebuild after code changes:
```bash
sudo docker compose build --no-cache
sudo docker compose up -d
```

Stop/Start:
```bash
sudo docker compose down
sudo docker compose up -d
```

Run Django admin commands:
```bash
sudo docker compose exec backend python manage.py createsuperuser
```

Migrations run automatically at container start, but you can run them manually as needed:
```bash
sudo docker compose exec backend python manage.py migrate
```

## 8) Troubleshooting
- 502/404 from `/server/`: check `frontend` logs and that `backend` is healthy and reachable as `http://backend:8000/` inside the network.
- CORS errors: when using the Nginx proxy path (`/server/`), CORS should not be needed. If calling the backend directly from a different origin, add that origin to `CORS_ALLOWED_ORIGINS` in Django settings.
- `ALLOWED_HOSTS` errors: add your domain/IP to Django `ALLOWED_HOSTS`.
