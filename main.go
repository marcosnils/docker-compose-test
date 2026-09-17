package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func redisAddr() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "redis"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

// getHitCount increments and returns the visit counter, retrying on
// connection errors to tolerate a slow-starting redis backend.
func getHitCount() (int64, error) {
	retries := 5
	for {
		count, err := rdb.Incr(ctx, "hits").Result()
		if err == nil {
			return count, nil
		}
		if retries == 0 {
			return 0, err
		}
		retries--
		time.Sleep(500 * time.Millisecond)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	count, err := getHitCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello World! I have been seen %d times.\n", count)
}

func main() {
	rdb = redis.NewClient(&redis.Options{Addr: redisAddr()})

	http.HandleFunc("/", hello)

	addr := ":5000"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
