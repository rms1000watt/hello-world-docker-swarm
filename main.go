package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var redisConn redis.Conn
var pgConn *sqlx.DB

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);
`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func main() {
	fmt.Println("Starting Server")

	if err := connect(); err != nil {
		fmt.Println("Failed connecting:", err)
		return
	}

	listenPort := os.Getenv("GLP_LISTEN_PORT")
	addr := fmt.Sprintf(":%s", listenPort)

	http.HandleFunc("/redis", LogMiddleware(HandlerRedis))
	http.HandleFunc("/pg", LogMiddleware(HandlerPg))

	fmt.Println("Listening on ", addr)
	http.ListenAndServe(addr, nil)
}

func connect() (err error) {
	redisHost := os.Getenv("GLP_REDIS_HOST")
	redisPort := os.Getenv("GLP_REDIS_PORT")

	fmt.Println("Connecting to Redis")
	redisConn, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
	if err != nil {
		fmt.Println("Failed connecting to Redis:", err)
		return err
	}

	pgHost := os.Getenv("GLP_PG_HOST")
	pgPort := os.Getenv("GLP_PG_PORT")
	pgUser := os.Getenv("GLP_PG_USER")
	pgPass := os.Getenv("GLP_PG_PASS")
	pgDb := os.Getenv("GLP_PG_DB")

	fmt.Println("Connecting to Postgres")
	dataSourceName := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDb)
	pgConn, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		fmt.Println("Failed connecting to Postgres:", err)
		return err
	}

	if _, err = pgConn.Exec(schema); err != nil {
		fmt.Println("Failed creating table:", err)
		// Allow this error to happen because.. maybe the table is created already.
		// I just need this server to run
	}

	if _, err = pgConn.Exec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Ryan", "Smith", "rms1000watt@test.com"); err != nil {
		fmt.Println("Failed inserting into table:", err)
		// Allow this error to happen because.. I just need this server to run
	}

	return
}

func LogMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling: /%s\n", r.URL.Path[1:])
		fn(w, r)
	}
}

func HandlerRedis(w http.ResponseWriter, r *http.Request) {
	reply, err := redisConn.Do("info")
	if err != nil {
		fmt.Println("Failed getting info:", err)
		return
	}

	info := InterfaceToString(reply)

	fmt.Println("redis:", info)
	w.Write([]byte(info))
}

func HandlerPg(w http.ResponseWriter, r *http.Request) {
	persons := []Person{}
	pgConn.Select(&persons, "SELECT * FROM person;")

	personsString := InterfaceToString(persons)

	fmt.Println("pg:", personsString)
	w.Write([]byte(personsString))
}

func InterfaceToString(in interface{}) (out string) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		fmt.Println("Failed encoding:", err)
		return "(This is a terrible way to handle errors....): " + err.Error()
	}

	return buf.String()
}
