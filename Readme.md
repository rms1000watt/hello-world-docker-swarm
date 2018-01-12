# Hello World Docker Swarm

## Introduction

Basic project using Docker Swarm

## Contents

- [Build](#build)
- [Usage](#usage)
- [Local Test](#local-test)

## Build

```bash
./build.sh
docker push rms1000watt/golang-redis-pg:latest
```

## Usage

```bash
# Start 2 service boxes and 1 db box
docker-machine create --driver virtualbox svc-1 &
docker-machine create --driver virtualbox svc-2 &
docker-machine create --driver virtualbox db-1  &

# List nodes
docker-machine ls

# Create a Swarm Master
docker-machine ssh svc-1 "docker swarm init --listen-addr $(docker-machine ip svc-1) --advertise-addr $(docker-machine ip svc-1)"
export WORKER_TOKEN=$(docker-machine ssh svc-1 "docker swarm join-token worker -q")

# Join a 'svc' Node to the cluster
docker-machine ssh svc-2 "docker swarm join --token=${WORKER_TOKEN} --listen-addr $(docker-machine ip svc-2) --advertise-addr $(docker-machine ip svc-2) $(docker-machine ip svc-1)"

# Join a 'db' Node to the cluster
docker-machine ssh db-1 "docker swarm join --token=${WORKER_TOKEN} --listen-addr $(docker-machine ip db-1) --advertise-addr $(docker-machine ip db-1) $(docker-machine ip svc-1)"

# Configure yourself to talk with the master
eval $(docker-machine env svc-1)

# View all nodes
docker node ls

# Label nodes
docker node update --label-add svc=true svc-1
docker node update --label-add svc=true svc-2
docker node update --label-add db=true --label-add pg-master=true db-1

# Create a network
docker network create --driver=overlay test-net

# Deploy test stack to swarm
docker stack deploy -c docker-compose.yml test-stack

# View what was deployed
docker stack ps test-stack
docker service ls

# View logs from the api server
docker service logs test-stack_golang-redis-pg

# Hit the server through the proxy
curl -H Host:test-stack-golang-redis-pg.traefik http://$(docker-machine ip svc-1)/redis
curl -H Host:test-stack-golang-redis-pg.traefik http://$(docker-machine ip svc-2)/redis
curl -H Host:test-stack-golang-redis-pg.traefik http://$(docker-machine ip db-1)/redis

# If you update your /etc/hosts with `192.168.99.117 test-stack-golang-redis-pg.traefik` you can: 
# curl http://test-stack-golang-redis-pg.traefik/redis

# When you're all done, delete the stack and the VMs
docker stack rm test-stack
docker-machine rm svc-1 svc-2 db-1

# List volumes
docker volume ls

# Remove unused volumes
docker volume prune

# Unset the env vars so you're looking at the host docker
eval $(docker-machine env -u)
```

## Local Test

```bash
GLP_LISTEN_PORT=9990 GLP_REDIS_HOST=192.168.99.103 GLP_REDIS_PORT=6379 GLP_PG_HOST=192.168.99.103 GLP_PG_PORT=5432 GLP_PG_USER=postgres GLP_PG_PASS=password GLP_PG_DB=postgres go run main.go
```
