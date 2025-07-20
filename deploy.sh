#!/bin/bash

set -e # Exit the script if any command fails

SERVICES=("user")

NAMESPACE=quizchief

if [ -z "$1" ]; then
    echo "Usage: ./deploy.sh [prod]"
    exit 1
fi

ENV="$1"
if [[ "$ENV" != "prod" ]]; then
    echo "Error: Invalid environment: '$ENV'. Must be or 'prod'."
    exit 1
fi

for SERVICE in "${SERVICES[@]}"; do
    IMAGE_NAME="quizchief-$SERVICE-service"
    VERSION=$(git rev-parse --short HEAD)
    FULL_IMAGE_NAME="$IMAGE_NAME:$VERSION"
    HELM_PATH="./server/helm/$SERVICE"
    DOCKER_REPO="codymmoore97"
    HELM_ENV_VALUES="$HELM_PATH/values.$ENV.yaml"

    echo -e "Pulling $DOCKER_REPO/$FULL_IMAGE_NAME..."
    docker pull $DOCKER_REPO/$FULL_IMAGE_NAME

    echo -e "Deploying $SERVICE..."
    helm upgrade --install $SERVICE "$HELM_PATH" --namespace $NAMESPACE --create-namespace \
        --set image.repository="$DOCKER_REPO/$IMAGE_NAME" \
        --set image.tag="$VERSION" \
        -f "$HELM_ENV_VALUES"
    kubectl rollout restart deployment $SERVICE -n quizchief
done

