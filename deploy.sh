#!/bin/bash

set -e # Exit the script if any command fails

SERVICES=("user")

NAMESPACE=quizchief

for SERVICE in "${SERVICES[@]}"; do
    IMAGE_NAME="quizchief-$SERVICE"
    VERSION="latest"
    HELM_PATH="./server/helm/$SERVICE"
    DOCKER_REPO="codymmoore97"

    echo -e "Building docker image for $SERVICE..."
    docker build -f server/internal/$SERVICE/Dockerfile -t $IMAGE_NAME:$VERSION server
    docker tag $IMAGE_NAME:$VERSION $DOCKER_REPO/$IMAGE_NAME:$VERSION
    echo -e "Pushing docker image to remote repository for $SERVICE..."
    docker push $DOCKER_REPO/$IMAGE_NAME:$VERSION

    echo -e "Deploying $SERVICE..."
    helm upgrade --install $SERVICE "$HELM_PATH" --namespace $NAMESPACE --create-namespace \
        --set image.repository="$DOCKER_REPO/$IMAGE_NAME" \
        --set image.tag="$VERSION"
done

