FROM golang:1.13

MAINTAINER theotimeriviere@gmail.com

ADD bin/fizzbuzz-http /usr/local/bin/fizzbuzz-http
ADD bin/fizzbuzz-grpc /usr/local/bin/fizzbuzz-grpc