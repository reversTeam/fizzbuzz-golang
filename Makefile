install: env protogen
	go get ./...

env:
	source .env

protogen:
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:. \
	src/**/protobuf/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	src/**/protobuf/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	src/**/protobuf/*.proto

clean:
	rm src/**/protobuf/*.pb.go || true
	rm src/**/protobuf/*.pb.gw.go || true
	rm src/**/protobuf/*.swagger.json || true

run:
	go run gateway.go