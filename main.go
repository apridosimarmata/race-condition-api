package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	opt := &goredislib.Options{
		Addr: "localhost:6379",
	}

	redisClient := goredislib.NewClient(opt)
	// pool := goredis.NewPool(redisClient)
	// rs := redsync.New(pool)
	// mutexname := "my-global-mutex"
	// mutex := rs.NewMutex(mutexname)

	// GET / -> setiap kali diakses, mengurangi nilai dari my-counter sebanyak 1
	// GET /set

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// lock
		reqId := r.URL.Query().Get("myId")
		// if err := mutex.Lock(); err != nil {
		// 	w.WriteHeader(409)
		// 	w.Write([]byte(string("oops! could not decrease.")))
		// 	return
		// }

		fmt.Printf("Processing request: %s\n", reqId)

		// fmt.Printf("I am locking: %s Mutex: %v\n", reqId, mutex.Name())

		// retrieving my-counter value
		result := redisClient.Get(context.Background(), "my-counter")
		if err := result.Err(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		myCounter, err := strconv.Atoi(result.Val())
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// should return 400 if the value already zero or less
		if myCounter < 1 {
			w.WriteHeader(400)
			w.Write([]byte(string("oops! could not decrease.")))
			return
		}

		// decreasing my-counter
		decreaseResult := redisClient.Decr(context.Background(), "my-counter")
		if err := decreaseResult.Err(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// // Release the lock so other processes or threads can obtain a lock.
		// if ok, err := mutex.Unlock(); !ok || err != nil {
		// 	w.WriteHeader(500)
		// 	return
		// }

		// fmt.Printf("I am releasing: %s Mutex: %v\n", reqId, mutex)
		fmt.Printf("Done, request: %s", reqId)

		// return success
		w.WriteHeader(200)
		w.Write([]byte(string("success")))
	})

	r.Get("/set", func(w http.ResponseWriter, r *http.Request) {

		result := redisClient.Set(context.Background(), "my-counter", 2, time.Second*600)
		if err := result.Err(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(string("value set to 2")))
	})

	http.ListenAndServe(":3000", r)
}
