# fizzbuzz-golang

- [How to setup the project](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/setup.md)

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz".

The output would look like this:
```
"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".
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
	    - Import Dashboard: 8588, 6671
	- [ ] Alerting
	  - https://devopscube.com/alert-manager-kubernetes-guide/
	- [ ] Terraformed (for the best)

  - Easy to maintain by other developers:
	- [ ] Linter
	- [ ] CI/CD : https://circleci.com/docs/
	- [ ] Changelog (for the best) : https://github.com/git-chglog/git-chglog
	- [ ] Fonctionnal & Unit test

## Bonus
  - [x] Add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:
	- Accept no parameter
	- Return the parameters corresponding to the most used request, as well as the number of hits for this request
