FROM google/golang:latest
MAINTAINER verdverm@gmail.com

VOLUME ["/logdir"]

ADD . /gopath
WORKDIR /gopath

CMD go run main.go
