# Finest project out there

## Web app

This is a simple http node server with EJS templates.

### How do I run it locally ?

```bash
cd webapp

# actually runs deno run --allow-net --allow-read --watch main.ts
# this allows the process to access network and reading fils such as the html template
deno task dev
```

And you're done ! The application is running on port 8080.

## Pushing images to docker

```bash
cd webapp

# build it
docker build -t ghcr.io/do3-2023/mmo-kube/webapp:<tag>

# test it
docker run -itp 8080:8080 ghcr.io/do3-2023/mmo-kube/webapp:<tag>

# push it
docker push ghcr.io/do3-2023/mmo-kube/webapp:<tag>
```

## Setting up k3d and the applications

```bash
mkdir /tmp/k3dvol

# -p options maps the nodeport 30080 to the 1024 port on your machine
# --volume option creates a volume on your machine so that persistent volumes can be retained
# --agents add two workers
# --api-port is the port of the kubeapi
k3d cluster create mmo --api-port 6550 -p "1024:30080@loadbalancer" --agents 2 --volume /tmp/k3dvol:/k3dvol

# Save the kubeconfig file, you can also put it somewhere else and specify KUBECONFIG env var
k3d kubeconfig get mmo > ~/.kube/config

# Wait for cluster to start
kubectl apply -f .kube/common/namespace

kubectl apply -f .kube/db

kubectl apply -f .kube/api

kubectl apply -f .kube/webapp

# wait for all deployments to be ready
curl http://localhost:1024

kubectl delete deployment -n data db

# theses pods not be ready and kubernetes should be trying to start them again
# the data is also persisted
kubectl get pods -n back
kubectl get pods -n front
```
