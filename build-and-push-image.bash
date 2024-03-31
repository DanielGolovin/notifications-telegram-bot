#!/bin/bash

# Load the environment variables from the .env file
set -a
source .env
set +a

docker build -f Dockerfile.prod -t $DOCKER_TAG .
docker tag $DOCKER_TAG $DOCKER_USERNAME/$DOCKER_REPO:$DOCKER_TAG
docker push $DOCKER_USERNAME/$DOCKER_REPO:$DOCKER_TAG
