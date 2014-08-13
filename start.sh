#!/bin/bash

# set -e

curr_dir="$(pwd)/$(dirname $0)"
echo "curr_dir = $curr_dir"

mkdir -p $curr_dir/logs

echo "Starting server httopd!"
docker run -d --name httopd-server \
	-p 8080:8080 \
	-p 8081:8081 \
	-p 8082:8082 \
	-v $curr_dir/logs:/var/log/nginx \
	verdverm/httopd-server > /dev/null

echo "Starting client httopd!"
docker run -d --name httopd-client \
	--net host \
	verdverm/httopd-client > /dev/null

# echo "Starting monitor httopd!"
# docker run -i -t --name httopd-monitor \
# 	verdverm/httopd-monitor
