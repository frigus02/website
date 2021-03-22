#!/bin/bash
set -eu

IMAGE=frigus02/website
if [ "$(git diff --stat)" != "" ]; then
    TAG="dev"
else
    TAG=$(git rev-parse HEAD)
fi

docker build -t "$IMAGE:$TAG" .
docker tag "$IMAGE:$TAG" "$IMAGE:latest"

if [ "$TAG" != "dev" ]; then
	if [ -z ${CI_REF+x} ] || [ "${CI_REF-x}" == "refs/heads/main" ]; then
		docker login -u frigus02 -p "$DOCKER_PASSWORD"
		docker push "$IMAGE:$TAG"
		docker push "$IMAGE:latest"
	fi
fi
