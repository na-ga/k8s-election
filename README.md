# k8s-election

Sample code for k8s leader election using ConfigMaps.

## Using minikube

Running

```
$ eval $(minikube docker-env)
$ make docker-build
$ kubectl run k8s-elector --image=k8s-elector:1.0 --replicas=3
```

Pod Status

```
$ kubectl get pod
NAME                           READY     STATUS    RESTARTS   AGE
k8s-elector-7fbd5b6869-7qcdm   1/1       Running   0          4m
k8s-elector-7fbd5b6869-dwcqk   1/1       Running   0          4m
k8s-elector-7fbd5b6869-xdjsj   1/1       Running   0          4m
```

Leader Election

```
$ k8stail -l run=k8s-elector
Context:   minikube
Namespace: default
Labels:    run=k8s-elector
Press Ctrl-C to exit.
----------
Pod:k8s-elector-7fbd5b6869-7qcdm Container:k8s-elector has been detected
Pod:k8s-elector-7fbd5b6869-xdjsj Container:k8s-elector has been detected
Pod:k8s-elector-7fbd5b6869-dwcqk Container:k8s-elector has been detected
[k8s-elector-7fbd5b6869-xdjsj][k8s-elector]  | 2017/12/25 08:45:41 I am a leader
[k8s-elector-7fbd5b6869-xdjsj][k8s-elector]  | 2017/12/25 08:45:46 I am a leader

...

[k8s-elector-7fbd5b6869-xdjsj][k8s-elector]  | 2017/12/25 08:46:11 I am a leader
[k8s-elector-7fbd5b6869-xdjsj][k8s-elector]  | 2017/12/25 08:46:13 Shutdown signal is received
[k8s-elector-7fbd5b6869-77xdh][k8s-elector]  | 2017/12/25 08:46:14 Detected leader: readerId=k8s-elector-7fbd5b6869-xdjsj
[k8s-elector-7fbd5b6869-7qcdm][k8s-elector]  | 2017/12/25 08:46:20 Detected leader: readerId=k8s-elector-7fbd5b6869-7qcdm
[k8s-elector-7fbd5b6869-7qcdm][k8s-elector]  | 2017/12/25 08:46:20 Started leading
[k8s-elector-7fbd5b6869-7qcdm][k8s-elector]  | 2017/12/25 08:46:20 Event(v1.ObjectReference{Kind:"ConfigMap", Namespace:"default", Name:"k8s-election", UID:"1e280461-e929-11e7-95bf-0800274ba730", APIVersion:"v1", ResourceVersion:"188945", FieldPath:""}): type: 'Normal' reason: 'LeaderElection' k8s-elector-7fbd5b6869-7qcdm became leader
[k8s-elector-7fbd5b6869-dwcqk][k8s-elector]  | 2017/12/25 08:46:20 Detected leader: readerId=k8s-elector-7fbd5b6869-7qcdm
[k8s-elector-7fbd5b6869-77xdh][k8s-elector]  | 2017/12/25 08:46:21 Detected leader: readerId=k8s-elector-7fbd5b6869-7qcdm
[k8s-elector-7fbd5b6869-7qcdm][k8s-elector]  | 2017/12/25 08:46:23 I am a leader
[k8s-elector-7fbd5b6869-7qcdm][k8s-elector]  | 2017/12/25 08:46:28 I am a leader
```