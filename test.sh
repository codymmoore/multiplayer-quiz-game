#!/bin/bash
set -euo pipefail

PROJECT_ROOT="$(pwd)"
SERVICES=("user")

echo "Running unit tests..."
for SERVICE in "${SERVICES[@]}"; do
    SERVICE_PATH="./server/internal/$SERVICE"

    cd $SERVICE_PATH

    echo "---------- ${SERVICES} ----------"
    go test -v

    cd $PROJECT_ROOT
done

echo -e "\nBuilding and starting docker containers..."
cd server
docker compose up --build -d
cd $PROJECT_ROOT

echo -e "\nRunning Postman tests..."
for SERVICE in "${SERVICES[@]}"; do
    echo "---------- ${SERVICE} ----------"
    newman run "postman/collection/Quizchief - ${SERVICE^} Service.postman_collection.json" \
        --environment "postman/environment/Quizchief - Local.postman_environment.json" \
        --reporters cli,json
done