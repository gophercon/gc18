## Gophercon 2018

This repository is a *gross* over-engineering of the 2018 Gophercon website to serve as a learning exercise, and to demonstrate some technology.

Components:

* Kubernetes
* Jaeger
* go-micro
* buffalo

## Installation/Setup

Step 0: Create a resource group in Azure.

    Mine is called "gophercon"

Step 1:  Setup Kubernetes Cluster
https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-kubernetes-walkthrough

```shell
    init/setupACS.sh
```

Step 2: Install jaeger 

```
    kubectl apply -f jaeger.yaml
    kubectl apply -f jaeger-ingress.yaml
```

Step 3: Create persistent volume for Consul/Traefik

```
    kubectl apply -f kubernetes/pv/storage-class.yaml
    kubectl apply -f kubernetes/pv/claim.yaml
```

Step 4: Install traefik & dependencies


```
    kubectl apply -f kubernetes/consul/*
```

Step 4: Build the web app 

```
    cd gophercon && docker build -t bketelsen/gc18-gophercon:<tag>
    docker push bketelsen/gc18-gophercon'
```

Step 5: deploy website

```
    kubectl apply -f kubernetes/website/*
```
    
