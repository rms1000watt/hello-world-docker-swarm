# Hello World Golang Redis

## Introduction

Basic project using Golang + Redis

## Contents

- [Build](#build)
- [Usage](#usage)

## Build

```bash
./build.sh
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

# Deploy to swarm
docker stack deploy -c docker-compose.yml test-stack

# View what was deployed
docker stack ps test-stack
docker service ls

# View logs from the api server
docker service logs test-stack_hello-world-golang-redis

# Try and hit the server
# curl http://192.168.99.105:9998/info
curl "http://$(docker-machine ls --filter name=svc-1 -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"
curl "http://$(docker-machine ls --filter name=svc-2 -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"
curl "http://$(docker-machine ls --filter name=db-1  -f '{{.URL}}' | awk -F'[/:]' '{printf $4}'):9998/info"

# When you're all done, delete the stack
docker stack rm test-stack
```
