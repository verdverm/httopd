#!/bin/bash

echo "stopping httopd-*!"
echo "----------------"

docker rm -f httopd-client
docker rm -f httopd-server

rm -rf $(dirname $0)/logs
