set -e

# docker related... to get the gateway IP
export HOSTURL=$(netstat -nr | grep 'UG' | awk '{print $2}')
echo "HOSTURL:   $HOSTURL"

go run main.go
