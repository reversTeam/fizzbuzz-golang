#!make


install: protogen
	go get ./...
	kubectl create namespace monitoring
	# kube prometheus
	kubectl create -f kubernetes/prometheus/clusterRole.yaml
	kubectl create -f kubernetes/prometheus/config-map.yaml
	kubectl create -f kubernetes/prometheus/prometheus-deployment.yaml
	kubectl create -f kubernetes/prometheus/prometheus-service.yaml
	# kube state metrics
	kubectl apply -f kubernetes/kube-state-metrics-configs/

	# kube grafana
	kubectl create -f kubernetes/grafana/datasource-config.yaml
	kubectl create -f kubernetes/grafana/deployment.yaml
	kubectl create -f kubernetes/grafana/service.yaml

	# redis
	kubectl create -f kubernetes/redis

linkports:
	kubectl port-forward service/grafana 3000:3000 -n monitoring > /dev/null &
	kubectl port-forward service/prometheus-service 9090:8080 -n monitoring > /dev/null &
	kubectl port-forward service/redis-master 6379:6379 > /dev/null &

protogen:
	for proto in src/endpoint/**/protobuf/*.proto ; do \
		protoc -I/usr/local/include -I. \
		  -I${GOPATH}/src \
		  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		  --go_out=plugins=grpc:. \
		$$proto ; \
		protoc -I/usr/local/include -I. \
		  -I${GOPATH}/src \
		  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		  --grpc-gateway_out=logtostderr=true:. \
		$$proto ; \
		protoc -I/usr/local/include -I. \
		  -I${GOPATH}/src \
		  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		  --swagger_out=logtostderr=true:. \
		$$proto ; \
	done

clean:
	rm src/endpoint/**/protobuf/*.pb.go || true
	rm src/endpoint/**/protobuf/*.pb.gw.go || true
	rm src/endpoint/**/protobuf/*.swagger.json || true

run:
	go run gateway.go

build:
	GOOS=linux GOARCH=amd64 go build -o fizzbuzz-http ./main.go
	GOOS=linux GOARCH=amd64 go build -o fizzbuzz-grpc ./src/endpoint/main.go

push:
	docker build -t triviere42/fizzbuzz-golang .
	docker push triviere42/fizzbuzz-golang

destroy:
	kubectl delete deployment client gateway || true
	kubectl delete service gateway client || true

apply:
	kubectl apply -f kubernetes/deployment.yaml
