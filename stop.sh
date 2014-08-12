#!/bin/bash

echo "stopping dawg-*!"
echo "----------------"

# docker rm -f dawg-nginx
docker rm -f dawg-server

rm -rf $curr_dir/logs
