#!/bin/bash

set -e

printf "\n\n\nbuilding server\n---------------------------\n"
cd server
docker build -t verdverm/dawg-server .
cd ..

printf "\n\n\nbuilding client\n---------------------------\n"
cd client
docker build -t verdverm/dawg-client .
cd ..

# printf "\n\n\nbuilding monitor\n---------------------------\n"
# cd monitor
# docker build -t verdverm/dawg-monitor .
# cd ..
