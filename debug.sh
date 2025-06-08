#!/bin/bash

set -e # Exit the script if any command fails

if [ -z "$1" ]; then
    echo "Usage: ./debug.sh <service> [--port=PORT]"
    exit 1
fi
SERVICE=$1

PORT=2345
for arg in "$@"; do
    if [[ "$arg" == --port=* ]]; then
        PORT="${arg#*=}"
        break
    fi
done

NAMESPACE="quizchief"
ENV="local"

IMAGE_NAME="quizchief-$SERVICE-service"
VERSION=$(git rev-parse --short HEAD)-debug
FULL_IMAGE_NAME="$IMAGE_NAME:$VERSION"
DOCKER_REPO="codymmoore97"

SERVICE_PATH="./server/internal/$SERVICE"

HELM_PATH="./server/helm/$SERVICE"
HELM_ENV_VALUES="$HELM_PATH/values.$ENV.yaml"

echo -e "Generating sqlc files for $SERVICE..."
sqlc generate --file $SERVICE_PATH/sqlc.yaml

echo -e "Building docker image for $SERVICE..."
docker build --no-cache --build-arg DEBUG=true --target=debug -f $SERVICE_PATH/Dockerfile -t $FULL_IMAGE_NAME server
docker tag $FULL_IMAGE_NAME $DOCKER_REPO/$FULL_IMAGE_NAME

echo -e "Pushing $DOCKER_REPO/$FULL_IMAGE_NAME..."
docker push $DOCKER_REPO/$FULL_IMAGE_NAME

echo -e "Deploying $SERVICE..."
helm upgrade --install $SERVICE "$HELM_PATH" --namespace $NAMESPACE --create-namespace \
    --set image.repository="$DOCKER_REPO/$IMAGE_NAME" \
    --set image.tag="$VERSION" \
    -f "$HELM_ENV_VALUES"
kubectl rollout restart deployment $SERVICE -n quizchief

POD=$(kubectl get pod -n $NAMESPACE -l app.kubernetes.io/name=$SERVICE --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1].metadata.name}')

echo "Waiting for pod $POD..."
kubectl wait -n $NAMESPACE --for=condition=ready pod $POD --timeout=60s

echo "Debugger listening on $PORT..."
kubectl port-forward pod/$POD $PORT:$PORT -n $NAMESPACE