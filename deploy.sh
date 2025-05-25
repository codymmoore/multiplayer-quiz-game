#!/bin/bash

set -e # Exit the script if any command fails

SERVICES=("user")

NAMESPACE=quizchief

if [ -z "$1" ]; then
    echo "Error: Environment not specified. Usage: ./deploy.sh [local|prod]"
    exit 1
fi

ENV="$1"
if [[ "$ENV" != "local" && "$ENV" != "prod" ]]; then
    echo "Error: Invalid environment '$ENV'. Must be 'local' or 'prod'."
    exit 1
fi

for SERVICE in "${SERVICES[@]}"; do
    IMAGE_NAME="quizchief-$SERVICE"
    VERSION="latest"
    FULL_IMAGE_NAME="$IMAGE_NAME:$VERSION"
    HELM_PATH="./server/helm/$SERVICE"
    DOCKER_REPO="codymmoore97"
    SERVICE_PATH="./server/internal/$SERVICE"
    HELM_ENV_VALUES="$HELM_PATH/values.$ENV.yaml"

    echo -e "Generating sqlc files for $SERVICE..."
    sqlc generate --file $SERVICE_PATH/sqlc.yaml

    echo -e "Building docker image for $SERVICE..."
    docker build -f $SERVICE_PATH/Dockerfile -t $FULL_IMAGE_NAME server
    docker tag $FULL_IMAGE_NAME $DOCKER_REPO/$FULL_IMAGE_NAME
    echo -e "Pushing docker image to remote repository for $SERVICE..."
    docker push $DOCKER_REPO/$FULL_IMAGE_NAME

    echo -e "Deploying $SERVICE..."
    helm upgrade --install $SERVICE "$HELM_PATH" --namespace $NAMESPACE --create-namespace \
        --set image.repository="$DOCKER_REPO/$FULL_IMAGE_NAME" \
        --set image.tag="$VERSION" \
        -f "$HELM_ENV_VALUES"
done

