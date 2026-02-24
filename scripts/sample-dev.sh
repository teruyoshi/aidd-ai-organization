#!/bin/bash

# AIDD TODO Sample - 開発環境起動スクリプト

set -e

COMPOSE_FILE="sample-docker-compose.yml"

echo "Starting AIDD TODO Sample development environment..."
docker-compose -f "$COMPOSE_FILE" up -d --build

echo ""
echo "Development environment started!"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  MySQL:    localhost:3306"
echo ""
echo "Logs: docker-compose -f $COMPOSE_FILE logs -f"
echo "Stop: docker-compose -f $COMPOSE_FILE down"
