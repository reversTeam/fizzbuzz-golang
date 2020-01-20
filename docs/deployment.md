# Fizzbuzz Golang Deployment
Currently the grpc service and the http over GRPC gateway are deployed on the same pod. This solution is sufficient in the case where we do not need a lot of computing power (~ 1.7k / req / node, ref standard-n1 GCP).

![Actual Deployment](https://raw.github.com/reversTeam/fizzbuzz-golang/master/assets/dashboard.jpg)

The pods are in rolling upgrade, that is to say that it will make a progressive rise in version during the pod update by checking that the new version works correctly, otherwise it will not go into production.

In this configuration it is impossible to communicate directly with the GRPC server in order to limit the infrastructure costs. Indeed, we note a CPU consumption on the part of the http gateway 2 times higher than the GRPC service (ratio observed on calls with a `limit` of 100).

![Actual Deployment](https://raw.github.com/reversTeam/fizzbuzz-golang/master/assets/deployment-v0.jpg)
```
apiVersion: v1
kind: Service
metadata:
  name: fizzbuzz
  labels:
    app: fizzbuzz
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/port: '4242'
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: fb-http-port
  selector:
    app: fizzbuzz
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fizzbuzz
spec:
  replicas: 6
  selector:
    matchLabels:
      app: fizzbuzz
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: fizzbuzz
    spec:
      containers:
        - name: http
          image: triviere42/fizzbuzz-golang
          command: ["/bin/sh"]
          args: ["-c", "fizzbuzz-http --http-host 0.0.0.0 --exporter-host=0.0.0.0"]
          ports:
            - containerPort: 8080
              name: fb-http-port
        - name: grpc
          image: triviere42/fizzbuzz-golang
          command: ["/bin/sh"]
          args: ["-c", "fizzbuzz-grpc --redis-host redis-master.default.svc.cluster.local"]
          ports:
            - containerPort: 42001
              name: fb-grpc-port
```

## How to improve design

We could have a better cost / scalability ratio. If we isolate the services on different pods with an exposure layer of the GRPC server via a loadbalancer. Then the HTTP gateways would connect randomly to the GRPC services.
The concern is that the connection between the server which acts as an http gateway and the GRPC service remains connected and that it is only possible to connect one service per path.
It would eventually be necessary to remake a component that would manage a pool of connection to GRPC services for an HTTP endpoint. This would allow us to have a better fault area.

It is also possible to connect other services than fizzbuzz, it is also possible to plan to create GRPC servers which would run on different pods and clusters.
If the need has been justified, the code allows it.

![Better Deployment](https://raw.github.com/reversTeam/fizzbuzz-golang/master/assets/deployment-v1.jpg)
```
apiVersion: v1
kind: Service
metadata:
  name: fizzbuzz-gw
  labels:
    app: fizzbuzz
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/port: '4242'
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: fb-http-port
  selector:
    app: fizzbuzz
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fizzbuzz-gw
spec:
  replicas: 6
  selector:
    matchLabels:
      app: fizzbuzz
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: fizzbuzz
    spec:
      containers:
        - name: http
          image: triviere42/fizzbuzz-golang
          command: ["/bin/sh"]
          args: ["-c", "fizzbuzz-http --http-host 0.0.0.0 --exporter-host=0.0.0.0 --grpc-host="fizzbuzz-grpc.default.svc.cluster.local"]
          ports:
            - containerPort: 8080
              name: fb-http-port
---
apiVersion: v1
kind: Service
metadata:
  name: fizzbuzz-grpc
  labels:
    app: fizzbuzz
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/port: '4242'
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: fb-grpc-port
  selector:
    app: fizzbuzz
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fizzbuzz-grpc
spec:
  replicas: 6
  selector:
    matchLabels:
      app: fizzbuzz
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: fizzbuzz
    spec:
      containers:
        - name: grpc
          image: triviere42/fizzbuzz-golang
          command: ["/bin/sh"]
          args: ["-c", "fizzbuzz-grpc --grpc-host 0.0.0.0 --redis-host redis-master.default.svc.cluster.local"]
          ports:
            - containerPort: 42001
              name: fb-grpc-port
```