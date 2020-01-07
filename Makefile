install: env protogen
	go get ./...

env:
	source .env

protogen:
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:. \
	protobuf/**/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	protobuf/**/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	protobuf/**/*.proto

clean:
	rm protobuf/**/*.pb.go || true
	rm protobuf/**/*.pb.gw.go || true
	rm protobuf/**/*.swagger.json || true

run:
	go run gateway.go