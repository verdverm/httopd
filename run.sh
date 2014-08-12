#!/bin/bash

# set -e

curr_dir="$(pwd)/$(dirname $0)"
echo "curr_dir = $curr_dir"

mkdir -p $curr_dir/logs

echo "Starting server dawg!"
docker run -d --name dawg-server \
	-p 5000:5000 \
	-p 8080:8080 \
	-v $curr_dir/logs:/var/log/nginx \
	verdverm/dawg-server > /dev/null

echo "Starting client dawg!"
docker run -d --name dawg-client \
	--net host \
	verdverm/dawg-client > /dev/null

# docker run -i -t --name dawg-client \
# 	--net host \
# 	verdverm/dawg-client

# echo "Starting monitor dawg!"
# docker run -i -t --name dawg-monitor \
# 	verdverm/dawg-monitor
