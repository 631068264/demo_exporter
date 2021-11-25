#!/bin/bash

set -ex

HUB_REPO="10.19.64.203:8080"
IMAGE_NAME=${HUB_REPO}/exporter/demo_exporter
WORK_DIR=/exporter

docker build --build-arg WORK_DIR=${WORK_DIR} -t ${IMAGE_NAME} -f Dockerfile .
docker run -it ${IMAGE_NAME} /bin/sh
#docker push ${IMAGE_NAME}

