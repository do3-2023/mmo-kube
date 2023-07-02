# Finest project out there

## Web app

This is a simple http node server with EJS templates.

Environment :

- `API_URL`: base url to the api

## Api

Simple golang http server using gin which connects to a postgresql database.

It also runs migrations at start if configured so.

Environment :

- `PG_USER`: Username of the user in the postgresql database
- `PG_PASSWORD`: Password of the user in the postgresql database
- `PG_HOSTNAME`: Hostname of the machine running postgresql
- `PG_DATABASE`: Name of the database to connect to
- `ENV`: Mode to run the api in. Possible values :
  - `dev`: runs migrations and start http server
  - `migrate`: runs migrations and exit

## How do I run it locally ?

#### Requirements

- [Docker Engine](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) as standalone binaries

#### Steps to follow

```bash
git clone https://github.com/do3-2023/mmo-kube

cd mmo-kube

# start the application
docker compose up -d

# access it
curl http://localhost:8080/
```

## Pushing images to a registry

```bash
cd webapp # or api

# build it
docker build -t "$IMAGE_NAME:$TAG"

# push it
docker push "$IMAGE_NAME:$TAG"
```

## Setting up k3d and the applications

### Requirements

- [Docker Engine](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [K3D](https://k3d.io/v5.5.1/)

### Steps to follow

#### Create this directory to persist the data later

```bash
mkdir /tmp/k3dvol
```

#### Create a cluster with the k3d cli

```bash
k3d cluster create mmo --api-port 6550 -p "1024:80@loadbalancer" --agents 2 --volume /tmp/k3dvol:/k3dvol
```

What does it do ?

- `-p` options maps the nodeport 80 which is the nodeport traefik is listening on to the 1024 port on your machine
- `--volume` option creates a volume on your machine so that persistent volumes can be retained
- `--agents` add two workers
- `--api-port` is the port of the kubeapi (optional)

#### Save the kubeconfig file, you can also put it somewhere else and specify KUBECONFIG env var

```bash
k3d kubeconfig get mmo > ~/.kube/config
```

#### Wait for cluster to start and apply the yaml files representing kubernetes resources

```bash
kubectl apply -f .kube/common/namespace

kubectl apply -f .kube/db

kubectl apply -f .kube/api

kubectl apply -f .kube/webapp
```

#### Wait for all deployments to be ready

```bash
curl http://localhost:1024
```

#### Delete the db deployment to check probes

```bash
kubectl delete deployment -n data db
```

#### Check the readiness of the webapp and api pods

```bash
kubectl get pods -n back
kubectl get pods -n front
```

They should not be ready.

#### Add the deployment again

```bash
kubectl apply -f .kube/db/deployment.yml
```

#### Check the readiness of the webapp and api pods

```bash
kubectl get pods -n back
kubectl get pods -n front
```

They should now be ready again. And the data should have been persisted.

### Cleanup

```bash
k3d cluster rm mmo

sudo rm -rf /tmp/k3dvol
```
