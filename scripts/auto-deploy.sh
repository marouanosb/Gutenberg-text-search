#!/usr/bin/env bash
set -euo pipefail

# Config (can be overridden via env)
REPO_DIR="${REPO_DIR:-$HOME/Gutenberg-text-search}"
REPO_URL="${REPO_URL:-https://github.com/marouanosb/Gutenberg-text-search.git}"

echo "[auto-deploy] REPO_DIR=${REPO_DIR}"
echo "[auto-deploy] REPO_URL=${REPO_URL}"

mkdir -p "${REPO_DIR}"
if [ ! -d "${REPO_DIR}/.git" ]; then
  echo "[auto-deploy] Cloning repo..."
  git clone "${REPO_URL}" "${REPO_DIR}"
fi

cd "${REPO_DIR}"
echo "[auto-deploy] Fetching latest..."
git fetch --all --prune
git reset --hard origin/main

# Choose docker or sudo docker
DOCKER="docker"
if ! ${DOCKER} ps >/dev/null 2>&1; then
  DOCKER="sudo docker"
fi

echo "[auto-deploy] Building and restarting compose stack..."
${DOCKER} compose pull || true
${DOCKER} compose build --pull
${DOCKER} compose up -d

echo "[auto-deploy] Pruning dangling images..."
${DOCKER} image prune -f || true

echo "[auto-deploy] Done."
