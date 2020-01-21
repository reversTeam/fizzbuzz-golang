# Why these choices?

### Why did you use GRPC?
GRPC is a communication protocol which has the particularity of being able to stay connected, moreover we can easily use definitions such as protobuf.
Protobuf allows us to centralize the definition of endpoints, we can then generate clients in different languages, the documentation of its endpoints, or even another layer of exposure such as GraphQl.


### Why use an http gateway?
It is a solution that seems to leave a lot of possibilities and integration. Could ultimately save on infrastructure by separating the http and grpc layer.
Protobuf defines the endpoints that will be exposed in HTTP with their parameters, the GRPC server faithfully follows the configuration of the HTTP endpoint.
So we have an HTTP server that will only accept parameters that are strictly defined in the prototype, and we will have a GRPC server that we are free to make public or not.
The resource consumption of a GRPC server is at least 2 times higher, according to the results I have observed in the realization of this project.

### Why did you use GCP as a provider for GCP?
Google is the company that initially created kubernetes, the idea was to quickly take over kubernetes in order to be able to check the scaling of the software.
I also benefit from free credit to discover the service. I strongly recommend that you put alerts on your billable consumption because it quickly becomes expensive.

### Why did you structure the project this way?
```
.
├── cmd                     // mains for generate bins
│   ├── grpc                // only for respect one main() by directory
│   └── http
└── src                     // The sources of project
    ├── common              // Common source between http, grpc or services
    └── endpoint            // The differents services
        └── fizzbuzz        // The fizzbuzz service
            └── protobuf    // Definition endpoint of fizzbuzz service
```
The idea was to make something small and easy to handle. It was not necessary that the code was cut out well not to be too oppressive, but also not too small not to have a nomenclature too weak to exploit.
The breakdown into services then seems natural, but I did not want the fact that supporting grpc and http to add a workload in the integration of the service.
So we had to standardize and reconcile the servers until we found common interfaces and a common boot process. The code that was, or that would be used for different components therefore needed to be centralized.

### Why did you set up dashboarding?
I did not really have visibility on the independent performance of my pods during my tests, I could possibly see a percentage of CPU usage of the pod which was running an instance of the program. The problem is that these figures were difficult to use. We must also look if two pods are not on the same node, or limited to deployment but that was not the objective.
It was therefore necessary to bring out the information in an exploitable format which allows me to see the figure that wanted on the loadbalancing of requests.
Know if when I down a pod the requests are well taken up by the other pods, or when I start a pod if it takes charge ...

### Why did you use a Loadblancer?
The GCP Loadbalancer was an easy-to-install layer that saves time, it does not correspond 100% as needed, or else I did not look at their documentation. But in the case of the gateway which connects to the GRPC server, the connection remains active, so when a gateway finds a grpc server through the LB it will remain connected to this server for the next requests.
To see what it is possible to do and if another LB corresponds better to the need, does not prevent mounted in chage.

### Why there is no HPA?
It is a small application which does not make of Io I did not have time to regader how I could set up metrics effectively to set up the HPA. In addition, the volume of requests managed by a pod supports a more than reasonable load.

### Why is there no competition to make fizzbuzz?
There is a branch that takes care of that the figures were not significant on a single-core instance, the operations are very simple for the CPU, the management of go-routines is not negligible. The gain is present but not significantly.
Logic applied in the code, launch go-routines by batch of 10k numbers. Well, but you have to supervise all of this by a maximum number of concurrent go-routines, otherwise we risk opening too many go-routines at once.
