# How to setup the project

This project is written in go and works with a set of dependencies. Its current version is tested and developed on Google Cloud Platform, because it notably uses annotations for LoadBalancing.
We cannot guarantee its functional state on another environment in the current state.

## Requirement
 - Have a Google Cloud Platform account
   - Kubernetes : v1.13.11-gke.14
   - Compute Engine: 3 nodes `n1-standard-1`
 - `gcloud`: SDK 245.0.0 see https://cloud.google.com/sdk/docs/quickstart-macos?hl=fr
  - Required python3
 - `kubectl`: v1.17.0 see https://kubernetes.io/fr/docs/tasks/tools/install-minikube/
 	- Linux see k3s, k8s, k9s
 - `go`: 1.13.5
 - `protobuff`: 3


## How to install

 1. Source the `.env`
Run the following command when you are in the project folder. It allows you to load the environmental variables necessary for the proper functioning of the project.
```
source .env
```

 2. Install Proto and GRPC tools
You will need to run the following commands in order to install all the tools necessary for compiling your proto files. These tools will allow you to centralize the definitions of your GRPC and HTTP endpoints in your proto files. You can also generate swagger documentation from the same definitions.
```
# Get the protobuf files
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

 3. Clone googleapis for the annotations.proto file (required for compile protobuf)
You need to clone the project on your $GOPATH/src/gihub.com
/!\ Maybe unusual, This step can be improved
For some reason that escapes me the command `go get github.com/googleapis/googleapis` does not work. In order not to have compilation errors on .proto files, you will have to add it manually.
```
mkdir -p $GOPATH/src/github.com/googleapis
git clone https://github.com/googleapis/googleapis $GOPATH/src/github.com/googleapis/
 ```

 4. Go to GCP, create a kubernetes cluster
Go to your GCP console to create your cluster, make a pool of 3 nodes of type `n1-standard-1`. Find out about the GCP areas to select the one closest to you. Latency has a strong impact on request competition (`europe-west3` is good choice for FR location).

 5. Connect your local to the GCP cluster
Connect to your GCP console, go to the `Compute > Kubernetes Engine > Clusters` section and click on `Connect` on your cluster, you should see a command that looks like this
```
gcloud container clusters get-credentials fizzbuzz-golang --zone europe-west3 --project fizzbuzz-golang
```

 6. Apply kubernetes annexe files
You can now launch the command which will allow you to install the additional services.
```
make install
```

 7. Run the fizzbuzz
You can now deploy the fizzbuzz with the following command
```
make apply
```

 8. Get the public IP
Wait one or two minutes for the service to instantiate the Loadbalancer and retrieve its public IP with the following command. If the External IP is `pending`, wait and retry later.
```
kubectl get service
NAME            TYPE           CLUSTER-IP   EXTERNAL-IP     PORT(S)           AGE
fizzbuzz-http   LoadBalancer   10.0.7.93    35.228.6.43     80:31409/TCP      156m
kubernetes      ClusterIP      10.0.0.1     <none>          443/TCP           7d2h
redis-master    ClusterIP      10.0.13.93   <none>          6379/TCP          40h
```

 9. Get the stats endpoint
We send a curl command to the stats endpoint, in order to verify its correct operation with empty data. The result should be this:

```
curl -X GET '35.228.75.112/fizzbuzz' | jq .
{
  "error": "No data found",
  "code": 2,
  "message": "No data found"
}
```

 10. Get fizzbuzz endpoint
Now we will send a fizzbuzz request then check the stats endpoint
```
curl -X POST 35.228.75.112/fizzbuzz -d '{"int1": 3, "int2": 5, "limit" : 101, "str1": "fizz"}
{"Items":["1","2","fizz","4","5","fizz","7","8","fizz","10","11","fizz","13","14","fizz","16","17","fizz","19","20","fizz","22","23","fizz","25","26","fizz","28","29","fizz","31","32","fizz","34","35","fizz","37","38","fizz","40","41","fizz","43","44","fizz","46","47","fizz","49","50","fizz","52","53","fizz","55","56","fizz","58","59","fizz","61","62","fizz","64","65","fizz","67","68","fizz","70","71","fizz","73","74","fizz","76","77","fizz","79","80","fizz","82","83","fizz","85","86","fizz","88","89","fizz","91","92","fizz","94","95","fizz","97","98","fizz","100","101"]}
```
```
curl -X GET '35.228.75.112/fizzbuzz' | jq .
{
  "Int1": 3,
  "Int2": 5,
  "Limit": 101,
  "Str1": "fizz",
  "Requests": "1"
}
```

## TIPS

 1. Update the Docker base image
You probably want to modify the deployment image to deploy other applications, you can modify it in the `kubernetes / deployment.yaml` file. You will also have the following command to build and push the docker directly on the registry.
```
make build
```

 2. Grafana / Prometheus / Redis
You may want to have access from your local machine on the cluster services, to do this you can run the following command:
```
make linkports
```
 - Grafana : http://127.0.0.1:3000/
 - Prometheus : http://127.0.0.1:9090/
 - Redis : http://127.0.0.1:6379/

 3. Siege
You can use `siege` to perform a light benchmark and get an idea of ​​the stability of the components. The concern is that the program does not allow having more than 250 clients simultaneously.
```
siege -c250 -t300S --content-type "application/json" 'http://35.228.75.112/fizzbuzz POST {"int1": 3, "int2": 5, "limit" : 10000, "str1": "fizz", "str2":"buzz"}'
```

 4. Vegeta
You can then use `vegeta` which has no limit on competition.
```
vegeta attack -duration=1200s -rate=6000 -targets=test/vegeta.list -output=/dev/null
```

You can also use it with jplot to have a graph of the requests, I invite you to look at the following link https://github.com/rs/jplot
