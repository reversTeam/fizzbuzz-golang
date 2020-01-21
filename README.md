# fizzbuzz-golang

![Architecture](https://raw.github.com/reversTeam/fizzbuzz-golang/master/docs/assets/fizzbuzz-architecture.jpg)

- [How to setup the project](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/setup.md)
- [How to use](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/how_to_use.md)
- [Deployment](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/deployment.md)
- [Architecture](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/architecture.md)

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz".

The output would look like this:
```
"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".
```

Example for get request with the following parameters:
 - int1: 3
 - int2: 5
 - limit: 20
 - str1: fizz
 - str2: buzz
```bash
curl -sX POST 127.0.0.1:8080/fizzbuzz -d '{"int1": 3, "int2": 5, "limit" : 20, "str1": "fizz", "str2":"buzz"}' | jq .
{"Items":["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16","17","fizz","19","buzz"]}
```

Example for get the most frequent requested (refer to bonus section)
```bash
curl -sX GET "127.0.0.1:8080/fizzbuzz" | jq .
{
  "Int1": "3",
  "Int2": "5",
  "Limit": "10",
  "Str1": "fizz",
  "Str2": "buzz",
  "Requests": "2347358"
}
```


## The goal
Implement a web server that will expose a REST API endpoint that: 
  - [x] Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
  - [x] Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.


## Checkpoint:
The server needs to be:
  - Ready for production:
	- [x] Kubernetes
	  - Develop on GCP
	- [x] LB & Availability
	  - Loadbalanced by GCP Loadbalancer, client and gateway is scalabled
	- [x] Monitoring: (cf. https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/)
	  - https://github.com/bibinwilson/kubernetes-prometheus.git
	  - https://github.com/devopscube/kube-state-metrics-configs.git
	  - https://devopscube.com/setup-grafana-kubernetes/
	    - Import Dashboard: [8588](https://grafana.com/grafana/dashboards/8588)
	    - Personnal metrics
	- [ ] Alerting
	  - https://devopscube.com/alert-manager-kubernetes-guide/
	- [ ] Terraformed (for the best)

  - Easy to maintain by other developers:
	- [x] Linter : https://github.com/golangci/golangci-lint
	- [x] CI : https://circleci.com/docs/
	  - [x] Fonctionnal 
	  - [ ] Unit test (for the best)
	- [ ] CD (for the best)
		- Required to add registry in GCP
		- Add variable in circle-ci for GCP access 
	- [ ] Changelog (for the best) : https://github.com/git-chglog/git-chglog

## Bonus
  - [x] Add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:
	- Accept no parameter
	- Return the parameters corresponding to the most used request, as well as the number of hits for this request
