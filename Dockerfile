ARG GO_DOCKER_VERSION

FROM golang:${GO_DOCKER_VERSION} AS builder
WORKDIR /go/src/github.com/rms1000watt/hello-world-golang-redis
COPY . .
# Do govendor install instead...
RUN go get \
  github.com/garyburd/redigo/redis \
  github.com/jmoiron/sqlx \
  github.com/lib/pq
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o bin/hello-world-golang-redis

FROM scratch
COPY --from=builder /go/src/github.com/rms1000watt/hello-world-golang-redis/bin/hello-world-golang-redis /hello-world-golang-redis
ENTRYPOINT [ "/hello-world-golang-redis" ]