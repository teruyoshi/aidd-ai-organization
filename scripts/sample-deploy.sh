#!/bin/bash

# AIDD TODO Sample - デプロイスクリプト

set -e

COMPOSE_FILE="sample-docker-compose.yml"

# 必須環境変数チェック
if [ -z "$MYSQL_ROOT_PASSWORD" ] || [ -z "$JWT_SECRET" ]; then
    echo "Error: MYSQL_ROOT_PASSWORD and JWT_SECRET environment variables are required"
    exit 1
fi

echo "Deploying AIDD TODO Sample..."
docker-compose -f "$COMPOSE_FILE" up -d --build

# ヘルスチェック
echo "Waiting for backend health check..."
for i in $(seq 1 30); do
    if curl -f -s "http://localhost:8080/health" > /dev/null 2>&1; then
        echo "Backend is healthy!"
        break
    fi
    sleep 3
done

echo ""
echo "Deployment complete!"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
