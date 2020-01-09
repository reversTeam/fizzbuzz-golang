FROM golang:1.13

MAINTAINER theotimeriviere@gmail.com

ADD gateway /usr/local/bin/fizzbuzz-gateway
ADD client /usr/local/bin/fizzbuzz-client

ENTRYPOINT ["fizzbuzz-gateway"]