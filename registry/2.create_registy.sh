#!/bin/sh
set -o errexit

# create registry container unless it already exists
reg_name='openmcp'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

# connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${reg_name}" || true



