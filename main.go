package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/llitfkitfk/GoHighPerformance/pkg/db"
	"github.com/llitfkitfk/GoHighPerformance/pkg/handler"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Our top-level router doesn't need CSRF protection: it's simple.

	// ... but our /api/* routes do, so we add it to the sub-router only.

	//r.Use(csrf.Protect([]byte("32-byte-long-auth-key")))

	conf, err := GetConfig()
	if err != nil {
		log.Printf("Error getting config [%s]", err)
		os.Exit(1)
	}

	var database db.DB
	switch conf.DBType {
	case "mem":
		database = db.NewMem()
	case "redis":
		redisOpts := &redis.Options{
			Addr:     conf.RedisHost,
			Password: conf.RedisPass,
			DB:       int(conf.RedisDB),
		}
		redisClient := redis.NewClient(redisOpts)
		database = db.NewRedis(redisClient)

	default:
		log.Printf("Error: no available DB type %s", conf.DBType)
		os.Exit(1)
	}

	router := mux.NewRouter()

	handler.NewCreateHandler(database).RegisterRoute(router)

	portStr := fmt.Sprintf(":%d", conf.Port)
	log.Printf("Serving on %s", portStr)
	log.Fatal(http.ListenAndServe(portStr, router))

}
