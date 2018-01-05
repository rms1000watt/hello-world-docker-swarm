package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
)

var redisConn redis.Conn

func main() {
	if err := connect(); err != nil {
		fmt.Println("Failed connecting:", err)
		return
	}

	fmt.Println("Starting Server")
	listenPort := os.Getenv("HWGR_LISTEN_PORT")

	http.HandleFunc("/info", HandlerInfo)
	http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
}

func connect() (err error) {
	redisHost := os.Getenv("HWGR_REDIS_HOST")
	redisPort := os.Getenv("HWGR_REDIS_PORT")

	fmt.Println("Connecting to Redis")
	redisConn, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
	if err != nil {
		fmt.Println("Failed connecting to Redis")
		return err
	}

	return
}

func HandlerInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling: /%s\n", r.URL.Path[1:])

	reply, err := redisConn.Do("info")
	if err != nil {
		fmt.Println("Failed getting info:", err)
		return
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(reply); err != nil {
		fmt.Println("Failed encoding:", err)
		return
	}

	fmt.Println("Redis Info:", buf.String())
	w.Write(buf.Bytes())
}
