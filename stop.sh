#!/bin/bash

echo "stopping dawg-*!"
echo "----------------"

docker rm -f dawg-client
docker rm -f dawg-server

rm -rf $(dirname $0)/logs
