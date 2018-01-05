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
printf "export MASTER_IP=%s\n" $(docker-machine ls --filter name=svc-1 -f '{{.URL}}' | awk -F'[/:]' '{printf $4}') > master-ip.sh
docker-machine scp master-ip.sh svc-1:/home/docker/master-ip.sh
docker-machine ssh svc-1 'eval $(cat master-ip.sh) && docker swarm init --advertise-addr $MASTER_IP' | grep "docker swarm join --token" > master-join.sh

# Join a 'svc' Node to the cluster
docker-machine scp master-join.sh svc-2:/home/docker/master-join.sh
docker-machine ssh svc-2 'sh master-join.sh'

# Join a 'db' Node to the cluster
docker-machine scp master-join.sh db-1:/home/docker/master-join.sh
docker-machine ssh db-1 'sh master-join.sh'

# View all nodes
docker-machine ssh svc-1 "docker node ls"

# Configure yourself to talk with the master
eval $(docker-machine env svc-1)

# Label nodes
docker node update --label-add svc=true svc-1
docker node update --label-add svc=true svc-2
docker node update --label-add db=true --label-add pg-master=true db-1

# Deploy to swarm
docker stack deploy -c docker-compose.yml test-stack

# View what was deployed
docker stack ps test-stack
docker service ls

# View logs from the api server
docker service logs test-stack_golang-redis-pg

# Try and hit the server
curl "http://$(docker-machine ls --filter name=svc-1 -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"
curl "http://$(docker-machine ls --filter name=svc-2 -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"
curl "http://$(docker-machine ls --filter name=db-1  -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"

# When you're all done, delete the stack and the VMs
docker stack rm test-stack
docker-machine rm svc-1 svc-2 db-1

# List volumes
docker volume ls

# Remove unused volumes
docker volume prune

# Sometimes volumes are in use with stopped containers.. docker rm those containers then docker volume prune
```

## Local Test

```bash
GLP_LISTEN_PORT=9990 GLP_REDIS_HOST=192.168.99.103 GLP_REDIS_PORT=6379 GLP_PG_HOST=192.168.99.103 GLP_PG_PORT=5432 GLP_PG_USER=postgres GLP_PG_PASS=password GLP_PG_DB=postgres go run main.go
```