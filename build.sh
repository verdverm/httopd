#!/bin/bash

set -e

echo "building server"
cd server
docker build -t verdverm/dawg-server .
cd ..

# echo "building nginx"
# cd nginx
# docker build --no-cache -t verdverm/dawg-nginx .
# cd ..

# echo "building client"
# cd client
# docker build -t verdverm/dawg-client .
# cd ..

# echo "building monitor"
# cd monitor
# docker build -t verdverm/dawg-monitor .
# cd ..
