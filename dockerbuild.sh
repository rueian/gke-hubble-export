#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

for bin in export; do
  docker build --rm -t rueian/gke-hubble-$bin:latest --target $bin .
done