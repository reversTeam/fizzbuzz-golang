FROM golang:1.13

MAINTAINER theotimeriviere@gmail.com

ADD fizzbuzz-http /usr/local/bin/fizzbuzz-http
ADD fizzbuzz-grpc /usr/local/bin/fizzbuzz-grpc