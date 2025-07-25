#!/bin/bash

set -e # Exit the script if any command fails

SERVICES=("user")

for SERVICE in "${SERVICES[@]}"; do
    IMAGE_NAME="quizchief-$SERVICE-service"
    VERSION=$(git rev-parse --short HEAD)
    FULL_IMAGE_NAME="$IMAGE_NAME:$VERSION"
    DOCKER_REPO="codymmoore97"
    SERVICE_PATH="./server/internal/$SERVICE"

    echo "Generating sqlc files for $SERVICE..."
    sqlc generate --file $SERVICE_PATH/sqlc.yaml

    echo "Building docker image for $SERVICE..."
    docker build --no-cache --target=release -f $SERVICE_PATH/Dockerfile -t $FULL_IMAGE_NAME server
    docker tag $FULL_IMAGE_NAME $DOCKER_REPO/$FULL_IMAGE_NAME
    echo "Pushing $DOCKER_REPO/$FULL_IMAGE_NAME..."
    docker push $DOCKER_REPO/$FULL_IMAGE_NAME
done