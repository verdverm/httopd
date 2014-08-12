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

# echo "Starting nginx dawg!"
# docker run -d --name dawg-nginx \
# 	-p 8080:8080 \
# 	-v $curr_dir/tmp/logs:/var/log/nginx \
# 	verdverm/dawg-nginx > /dev/null

# echo "Starting nginx dawg!"
# docker run -i -t --name dawg-nginx \
# 	verdverm/dawg-nginx /bin/bash



# echo "Starting client dawg!"
# docker run -d --name dawg-client \
# 	verdverm/dawg-client > dev/null

# echo "Starting monitor dawg!"
# docker run -i -t --name dawg-monitor \
# 	verdverm/dawg-monitor
