#!/bin/bash

set -e

# docker related... to get the gateway IP
export AHOST=$(netstat -nr | grep 'UG' | awk '{print $2}')
echo "AHOST:   $AHOST"

(
	echo "# Auto-generated : do not touch"
	echo ""
	# echo "upstream site.localhost {"
	# printf "    server ${AHOST}:5000;\n"
	# echo "}"
	echo "server {"
	echo ""
	echo "    listen 8080;"
	echo "    server_name httopd-1;"
	echo "    access_log  /var/log/nginx/host.httopd-1.access.log;"
	echo ""
	echo "    location / {"
	echo "        proxy_set_header        Host            \$host;"
	echo "        proxy_set_header        X-Real-IP       \$remote_addr;"
	echo "        proxy_set_header        X-Forwarded-For \$proxy_add_x_forwarded_for;"
	echo ""
	echo "        proxy_pass http://127.0.0.1:5001;"
	echo "    }"
	echo "}"
	echo "server {"
	echo ""
	echo "    listen 8081;"
	echo "    server_name httopd-2;"
	echo "    access_log  /var/log/nginx/host.httopd-2.access.log;"
	echo ""
	echo "    location / {"
	echo "        proxy_set_header        Host            \$host;"
	echo "        proxy_set_header        X-Real-IP       \$remote_addr;"
	echo "        proxy_set_header        X-Forwarded-For \$proxy_add_x_forwarded_for;"
	echo ""
	echo "        proxy_pass http://127.0.0.1:5002;"
	echo "    }"
	echo "}"
	echo "server {"
	echo ""
	echo "    listen 8082;"
	echo "    server_name httopd-3;"
	echo "    access_log  /var/log/nginx/host.httopd-3.access.log;"
	echo ""
	echo "    location / {"
	echo "        proxy_set_header        Host            \$host;"
	echo "        proxy_set_header        X-Real-IP       \$remote_addr;"
	echo "        proxy_set_header        X-Forwarded-For \$proxy_add_x_forwarded_for;"
	echo ""
	echo "        proxy_pass http://127.0.0.1:5003;"
	echo "    }"
	echo "}"
) > dlb.cfg

sudo cp dlb.cfg /etc/nginx/sites-enabled/default

nginx &

export PORT=5001
python app.py &

export PORT=5002
python app.py &

export PORT=5003
python app.py
