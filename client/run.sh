set -e

# docker related... to get the gateway IP
# export HOSTURL=$(netstat -nr | grep 'UG' | awk '{print $2}')
# echo "HOSTURL:   $HOSTURL"
# go build

echo "client 8080"
go run *.go -host="localhost:8080" &

echo "client 8081"
go run *.go -host="localhost:8081" &

echo "client 8082"
go run *.go -host="localhost:8082"

