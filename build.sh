#!/usr/bin/env bash

docker build \
  --build-arg GO_DOCKER_VERSION=${GO_DOCKER_VERSION:-1.9.2} \
  --force-rm \
  --no-cache \
  --compress \
  -t rms1000watt/golang-redis-pg:latest .
