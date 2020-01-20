## How to use?

You can retrieve the service ip by issuing the following command, the ip may be pending, in this case you will have to wait for the Loadbalancer to finish installing.

```bash
kubectl get service
NAME           TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)        AGE
fizzbuzz       LoadBalancer   10.0.14.200   35.228.75.112   80:30531/TCP   36h
kubernetes     ClusterIP      10.0.0.1      <none>          443/TCP        8d
redis-master   ClusterIP      10.0.13.93    <none>          6379/TCP       3d4h
```


Now that you have the service IP address (in our example: 35.228.75.112), you can now try the endpoint with the following commands.

```bash
curl -sX POST 35.228.75.112/fizzbuzz -d '{"int1": 3, "int2": 5, "limit" : 100, "str1": "fizz", "str2":"buzz"}'
{"Items":["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16","17","fizz","19","buzz","fizz","22","23","fizz","buzz","26","fizz","28","29","fizzbuzz","31","32","fizz","34","buzz","fizz","37","38","fizz","buzz","41","fizz","43","44","fizzbuzz","46","47","fizz","49","buzz","fizz","52","53","fizz","buzz","56","fizz","58","59","fizzbuzz","61","62","fizz","64","buzz","fizz","67","68","fizz","buzz","71","fizz","73","74","fizzbuzz","76","77","fizz","79","buzz","fizz","82","83","fizz","buzz","86","fizz","88","89","fizzbuzz","91","92","fizz","94","buzz","fizz","97","98","fizz","buzz","101"]}
```

You can also send a GET request without parameters to the same URL, this will then send you the information on the endpoint most consulted with the number of requests that were made with these parameters.

```bash
curl -sX GET "35.228.75.112/fizzbuzz"
{"Int1":3,"Int2":5,"Limit":10,"Str1":"fizz","Str2":"buzz","Requests":"12123191"}
```

How to access services such as prometheus, grafana or redis from your local machine.
You can do port forwarding with the following commands:
```bash
kubectl port-forward service/grafana 3000:3000 -n monitoring > /dev/null &
kubectl port-forward service/prometheus-service 9090:8080 -n monitoring > /dev/null &
kubectl port-forward service/redis-master 6379:6379 > /dev/null &
```

All of its commands will be executed if you start the command `make linkports`

### Grafana

![Actual Deployment](https://raw.github.com/reversTeam/fizzbuzz-golang/master/docs/assets/dashboard.png)

The first time you use grafana you will have to connect with the default logins:
 - User: admin
 - Pass: admin

Then you can configure the password you want.

I recommend you install the following [dashboard](https://grafana.com/grafana/dashboards/8588) to analyze the kubernetes cluster, the `kube-state-netrics-config` service deploys a service in the namespace` kube-dns` which allows you to discover services while you resize your cluser.

You can also create a new dashboard, call it for example `FizzBuzz` then prepare 3 graphs that you plug into the backend` prometheus`, then apply the following queries to them:

This query allows you to add up all the requests that are managed on all of the pods in the cluster
```
sum(rate(fizzbuzz_request_sec[$__interval]))
```

This query allows you to have the number of requests managed by pod, to monitor that the loadblancing works correctly for example.
```
sum by (instance) (rate(fizzbuzz_request_sec[$__interval]))
```

This query allows you to group the queries by return code, ideal for quickly identifying problems during production. Beware of false positive, if a customer generates a lot of Bad Request it will make an unusual, but not abnormal spike.
```
sum by (code) (rate(fizzbuzz_request_sec[$__interval]))
```