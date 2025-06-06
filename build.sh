#!/bin/bash

set -e # Exit the script if any command fails

SERVICES=("user")

for SERVICE in "${SERVICES[@]}"; do
    IMAGE_NAME="quizchief-$SERVICE"
    VERSION=$(git rev-parse --short HEAD)
    FULL_IMAGE_NAME="$IMAGE_NAME:$VERSION"
    DOCKER_REPO="codymmoore97"
    SERVICE_PATH="./server/internal/$SERVICE"

    echo -e "Generating sqlc files for $SERVICE..."
    sqlc generate --file $SERVICE_PATH/sqlc.yaml

    echo -e "Building docker image for $SERVICE..."
    docker build -f $SERVICE_PATH/Dockerfile -t $FULL_IMAGE_NAME server
    docker tag $FULL_IMAGE_NAME $DOCKER_REPO/$FULL_IMAGE_NAME
    echo -e "Pushing docker image to remote repository for $SERVICE..."
    docker push $DOCKER_REPO/$FULL_IMAGE_NAME
done